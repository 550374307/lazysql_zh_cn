package app

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/jorgerojas26/lazysql/internal/mylogin"
	"github.com/jorgerojas26/lazysql/models"
)

// writeMyloginFile builds a .mylogin.cnf-compatible file with the given
// plaintext lines, mirroring mysql_config_editor's layout (used only in tests).
func writeMyloginFile(t *testing.T, dir string, lines []string) string {
	t.Helper()
	fileKey := []byte("0123456789abcdefghij") // 20 bytes

	// Same deriveKey XOR-fold as internal/mylogin.
	aesKey := make([]byte, 16)
	copy(aesKey, fileKey[:16])
	for i := 16; i < len(fileKey); i++ {
		aesKey[i-16] ^= fileKey[i]
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		t.Fatalf("aes: %v", err)
	}
	var buf bytes.Buffer
	buf.Write([]byte{0, 0, 0, 0})
	buf.Write(fileKey)
	for _, line := range lines {
		pad := aes.BlockSize - len(line)%aes.BlockSize
		padded := append([]byte(line), bytes.Repeat([]byte{byte(pad)}, pad)...)
		cipher := make([]byte, len(padded))
		for i := 0; i < len(padded); i += aes.BlockSize {
			block.Encrypt(cipher[i:i+aes.BlockSize], padded[i:i+aes.BlockSize])
		}
		var size [4]byte
		binary.LittleEndian.PutUint32(size[:], uint32(len(cipher)))
		buf.Write(size[:])
		buf.Write(cipher)
	}

	path := filepath.Join(dir, ".mylogin.cnf")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestApplyLoginPathFrom(t *testing.T) {
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{
		"[prod]",
		`user = "root"`,
		`password = "s3cret"`,
		`host = "127.0.0.1"`,
		`port = 3306`,
	})

	conn := models.Connection{LoginPath: "prod", DBName: "app"}
	if err := applyLoginPathFrom(&conn, path); err != nil {
		t.Fatalf("applyLoginPathFrom: %v", err)
	}
	if conn.Username != "root" || conn.Password != "s3cret" ||
		conn.Hostname != "127.0.0.1" || conn.Port != "3306" {
		t.Fatalf("conn = %+v", conn)
	}

	// Explicit values always win.
	conn2 := models.Connection{LoginPath: "prod", Username: "app", Port: "3307"}
	if err := applyLoginPathFrom(&conn2, path); err != nil {
		t.Fatalf("applyLoginPathFrom: %v", err)
	}
	if conn2.Username != "app" || conn2.Port != "3307" || conn2.Password != "s3cret" {
		t.Fatalf("conn2 = %+v", conn2)
	}
}

func TestApplyLoginPathMissingEntry(t *testing.T) {
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{"[prod]", `user = "root"`})
	conn := models.Connection{LoginPath: "nope"}
	if err := applyLoginPathFrom(&conn, path); err == nil {
		t.Fatal("expected error for missing login path")
	}
}

func TestParseConfigURLMySQLWithLoginPath(t *testing.T) {
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{
		"[prod]",
		`user = "root"`,
		`password = "s3cret"`,
		`host = "127.0.0.1"`,
		`port = 3306`,
	})

	conn := models.Connection{Provider: "mysql", LoginPath: "prod", DBName: "app"}
	if err := applyLoginPathFrom(&conn, path); err != nil {
		t.Fatalf("applyLoginPathFrom: %v", err)
	}
	got := parseConfigURL(&conn)
	want := "mysql://root:s3cret@127.0.0.1:3306/app"
	if got != want {
		t.Fatalf("parseConfigURL = %q, want %q", got, want)
	}
}

func TestParseConfigURLMySQLDefaultPort(t *testing.T) {
	conn := models.Connection{
		Provider:  "mysql",
		Username:  "root",
		Password:  "pw",
		Hostname:  "db.internal",
		DBName:    "app",
		URLParams: "?charset=utf8mb4",
	}
	got := parseConfigURL(&conn)
	want := "mysql://root:pw@db.internal:3306/app?charset=utf8mb4"
	if got != want {
		t.Fatalf("parseConfigURL = %q, want %q", got, want)
	}
}

func TestMyloginDefaultsParse(t *testing.T) {
	// Sanity: the mylogin package can decode what we write (cross-check of
	// the wire format used by both writer and reader).
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{
		"[client]",
		`user = "root"`,
		`password = "pw"`,
	})
	f, err := mylogin.Load(path)
	if err != nil {
		t.Fatalf("mylogin.Load: %v", err)
	}
	entry, ok := f.Get("client")
	if !ok {
		t.Fatal("login path 'client' not found")
	}
	if entry.User != "root" || entry.Password != "pw" {
		t.Fatalf("entry = %+v", entry)
	}
}
