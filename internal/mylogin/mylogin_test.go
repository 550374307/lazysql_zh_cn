package mylogin

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"testing"
)

// encryptECB mirrors the file format writer (used only for the round-trip
// test): AES-128-ECB with PKCS#7 padding.
func encryptECB(plain, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	pad := aes.BlockSize - len(plain)%aes.BlockSize
	padded := append(plain, bytes.Repeat([]byte{byte(pad)}, pad)...)
	out := make([]byte, len(padded))
	for i := 0; i < len(padded); i += aes.BlockSize {
		block.Encrypt(out[i:i+aes.BlockSize], padded[i:i+aes.BlockSize])
	}
	return out, nil
}

// writeFile builds a .mylogin.cnf byte layout with the given 20-byte key.
func writeFile(t *testing.T, fileKey []byte, lines []string) []byte {
	t.Helper()
	if len(fileKey) != 20 {
		t.Fatalf("fileKey must be 20 bytes, got %d", len(fileKey))
	}
	var buf bytes.Buffer
	buf.Write([]byte{0, 0, 0, 0}) // flag
	buf.Write(fileKey)

	aesKey := deriveKey(fileKey)
	for _, line := range lines {
		cipher, err := encryptECB([]byte(line), aesKey)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}
		var size [4]byte
		binary.LittleEndian.PutUint32(size[:], uint32(len(cipher)))
		buf.Write(size[:])
		buf.Write(cipher)
	}
	return buf.Bytes()
}

func TestParseRoundTrip(t *testing.T) {
	fileKey := []byte("0123456789abcdefghij") // 20 bytes
	data := writeFile(t, fileKey, []string{
		"[client]",
		`user = "root"`,
		`password = "s3cret"`,
		`host = "127.0.0.1"`,
		`port = 3306`,
		"[other]",
		`user = "app"`,
		`password = "pw"`,
	})

	f, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	client, ok := f.Get("client")
	if !ok {
		t.Fatal("login path 'client' not found")
	}
	if client.User != "root" || client.Password != "s3cret" ||
		client.Host != "127.0.0.1" || client.Port != "3306" {
		t.Fatalf("client = %+v", client)
	}

	other, ok := f.Get("other")
	if !ok || other.User != "app" || other.Password != "pw" {
		t.Fatalf("other = %+v, ok=%v", other, ok)
	}

	if _, ok := f.Get("missing"); ok {
		t.Fatal("unexpected login path 'missing'")
	}
}

func TestParseDerivedKeyXor(t *testing.T) {
	// fileKey bytes 16..19 fold into bytes 0..3 via XOR.
	fileKey := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x21, 0x22, 0x23, 0x24,
	}
	got := deriveKey(fileKey)
	want := []byte{
		0x20, 0x20, 0x20, 0x20, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("deriveKey[%d] = %#x, want %#x", i, got[i], want[i])
		}
	}
}

func TestParseTooShort(t *testing.T) {
	if _, err := Parse([]byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for short file")
	}
}
