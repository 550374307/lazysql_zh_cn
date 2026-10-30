// Package mylogin reads MySQL's encrypted login path file (.mylogin.cnf)
// as written by mysql_config_editor, so LazySQL can use credentials stored
// outside the plaintext config file.
//
// File layout (little-endian lengths):
//
//	[0:4]   flag (unused padding)
//	[4:24]  20-byte random key stored in the file
//	loop:
//	  [4 bytes]  length of the next AES-128-ECB ciphertext block
//	  [length]   ciphertext, plaintext is one "key = \"value\"" line or a
//	            "[login-path]" section header line
//
// The real AES key is derived from the 20-byte file key by XOR-folding it
// into 16 bytes: aesKey[j] = fileKey[j] ^ fileKey[j+16] for j in 0..3 and
// aesKey[j] = fileKey[j] otherwise. This mirrors MySQL's my_aes_keygen /
// mysql_config_editor behaviour (AES-128-ECB + PKCS#7 padding).
package mylogin

import (
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DefaultPath returns the mylogin.cnf location for the current user,
// matching the MySQL client: %APPDATA%\MySQL\.mylogin.cnf on Windows and
// ~/.mylogin.cnf elsewhere.
func DefaultPath() (string, error) {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "MySQL", ".mylogin.cnf"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".mylogin.cnf"), nil
}

// LoginPath holds the credentials stored for one --login-path entry.
type LoginPath struct {
	User     string
	Password string
	Host     string
	Port     string
	Socket   string
}

// File represents a parsed .mylogin.cnf.
type File struct {
	// Paths maps a login-path name (e.g. "client") to its credentials.
	Paths map[string]LoginPath
}

// Load parses the mylogin.cnf at path. A missing file is an error.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse decodes the binary contents of a .mylogin.cnf.
func Parse(data []byte) (*File, error) {
	if len(data) < 24 {
		return nil, errors.New("mylogin.cnf: file too short")
	}

	// The first 4 bytes are a filler flag; the next 20 are the file key.
	fileKey := data[4:24]

	aesKey := deriveKey(fileKey)

	offset := 24
	f := &File{Paths: make(map[string]LoginPath)}
	current := ""
	for offset+4 <= len(data) {
		size := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4
		if size <= 0 || offset+size > len(data) {
			return nil, fmt.Errorf("mylogin.cnf: invalid record size %d at offset %d", size, offset-4)
		}
		plain, err := decryptECB(data[offset:offset+size], aesKey)
		if err != nil {
			return nil, fmt.Errorf("mylogin.cnf: decrypting record at offset %d: %w", offset-4, err)
		}
		offset += size

		line := strings.TrimSpace(string(plain))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = line[1 : len(line)-1]
			if _, ok := f.Paths[current]; !ok {
				f.Paths[current] = LoginPath{}
			}
			continue
		}
		if current == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		field := strings.ToLower(strings.TrimSpace(key))
		val := strings.Trim(strings.TrimSpace(value), `"`)
		switch field {
		case "user":
			entry := f.Paths[current]
			entry.User = val
			f.Paths[current] = entry
		case "password":
			entry := f.Paths[current]
			entry.Password = val
			f.Paths[current] = entry
		case "host":
			entry := f.Paths[current]
			entry.Host = val
			f.Paths[current] = entry
		case "port":
			entry := f.Paths[current]
			entry.Port = val
			f.Paths[current] = entry
		case "socket":
			entry := f.Paths[current]
			entry.Socket = val
			f.Paths[current] = entry
		}
	}
	return f, nil
}

// Get returns the credentials for the named login path.
func (f *File) Get(name string) (LoginPath, bool) {
	lp, ok := f.Paths[name]
	return lp, ok
}

// deriveKey folds the 20-byte stored key into the 16-byte AES key.
func deriveKey(fileKey []byte) []byte {
	aesKey := make([]byte, 16)
	for i := range 16 {
		aesKey[i] = fileKey[i]
	}
	for i := 16; i < len(fileKey); i++ {
		aesKey[i-16] ^= fileKey[i]
	}
	return aesKey
}

// decryptECB decrypts one AES-128-ECB block sequence and removes PKCS#7
// padding.
func decryptECB(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext length is not a multiple of block size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plain := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(plain[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}
	// PKCS#7 unpadding.
	if len(plain) == 0 {
		return nil, nil
	}
	pad := int(plain[len(plain)-1])
	if pad == 0 || pad > aes.BlockSize || pad > len(plain) {
		return nil, errors.New("invalid PKCS#7 padding")
	}
	return plain[:len(plain)-pad], nil
}
