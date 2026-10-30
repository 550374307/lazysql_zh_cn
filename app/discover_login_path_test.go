package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"

	"github.com/jorgerojas26/lazysql/models"
)

func TestDiscoverLoginPathConnections(t *testing.T) {
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{
		"[prod]",
		`user = "root"`,
		`password = "s3cret"`,
		`host = "127.0.0.1"`,
		`port = 3306`,
		"[dev]",
		`user = "app"`,
		`password = "devpw"`,
	})

	conns, err := discoverLoginPathConnections(nil, path)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(conns) != 2 {
		t.Fatalf("got %d connections, want 2", len(conns))
	}

	prod := findByLoginPath(t, conns, "prod")
	if prod.Name != "prod" || prod.Provider != "mysql" || prod.Username != "root" ||
		prod.Password != "s3cret" || prod.Hostname != "127.0.0.1" || prod.Port != "3306" {
		t.Fatalf("prod = %+v", prod)
	}
	if !prod.Ephemeral {
		t.Fatal("auto-discovered connection must be ephemeral")
	}
	if prod.URL != "mysql://root:s3cret@127.0.0.1:3306/" {
		t.Fatalf("prod URL = %q", prod.URL)
	}

	dev := findByLoginPath(t, conns, "dev")
	if dev.URL != "mysql://app:devpw@localhost:3306/" {
		t.Fatalf("dev URL = %q (host should default to localhost)", dev.URL)
	}
}

func findByLoginPath(t *testing.T, conns []models.Connection, name string) models.Connection {
	t.Helper()
	for _, c := range conns {
		if c.LoginPath == name {
			return c
		}
	}
	t.Fatalf("connection %q not found in %+v", name, conns)
	return models.Connection{}
}

func TestDiscoverLoginPathSkipsExplicit(t *testing.T) {
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{
		"[prod]",
		`user = "root"`,
		"[dev]",
		`user = "app"`,
	})

	explicit := []models.Connection{
		{Name: "My Prod", Provider: "mysql", LoginPath: "prod", Username: "explicit", Password: "pw"},
	}
	conns, err := discoverLoginPathConnections(explicit, path)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(conns) != 2 {
		t.Fatalf("got %d connections, want 2 (explicit prod + auto dev)", len(conns))
	}
	if conns[0].Name != "My Prod" || conns[0].Username != "explicit" {
		t.Fatalf("explicit row changed: %+v", conns[0])
	}
	if conns[1].Name != "dev" || conns[1].Ephemeral != true {
		t.Fatalf("auto row wrong: %+v", conns[1])
	}
}

func TestDiscoverLoginPathMissingFile(t *testing.T) {
	conns := []models.Connection{{Name: "only", Provider: "mysql", URL: "mysql://u:p@h:3306/db"}}
	got, err := discoverLoginPathConnections(conns, "/nonexistent/.mylogin.cnf")
	if err != nil {
		t.Fatalf("missing file must not error, got %v", err)
	}
	if len(got) != 1 || got[0].Name != "only" {
		t.Fatalf("got = %+v", got)
	}
}

func TestDiscoverLoginPathSocket(t *testing.T) {
	dir := t.TempDir()
	path := writeMyloginFile(t, dir, []string{
		"[local]",
		`user = "root"`,
		`password = "pw"`,
		`socket = "/var/run/mysqld/mysqld.sock"`,
	})

	conns, err := discoverLoginPathConnections(nil, path)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(conns) != 1 {
		t.Fatalf("got %d connections", len(conns))
	}
	want := "mysql://root:pw@unix(/var/run/mysqld/mysqld.sock)/"
	if conns[0].URL != want {
		t.Fatalf("URL = %q, want %q", conns[0].URL, want)
	}
}

func TestEphemeralNotPersisted(t *testing.T) {
	// The toml:"-" tag keeps the marker field itself out of serialized output.
	conn := models.Connection{
		Name:      "auto",
		Provider:  "mysql",
		Username:  "root",
		Password:  "pw",
		Hostname:  "localhost",
		LoginPath: "auto",
		Ephemeral: true,
	}
	out, err := toml.Marshal([]models.Connection{conn})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(out), "ephemeral") {
		t.Fatalf("ephemeral marker leaked into config: %s", out)
	}
}

func TestSaveConnectionsSkipsEphemeral(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte("[application]\nDefaultPageSize = 100\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	c := &Config{ConfigFile: cfgPath}

	auto := models.Connection{Name: "auto", Provider: "mysql", LoginPath: "auto", Ephemeral: true}
	real := models.Connection{Name: "real", Provider: "mysql", URL: "mysql://u:p@h:3306/db"}
	if err := c.SaveConnections([]models.Connection{auto, real}); err != nil {
		t.Fatalf("SaveConnections: %v", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if strings.Contains(string(data), "'auto'") {
		t.Fatalf("auto-discovered connection leaked into config file: %s", data)
	}
	if !strings.Contains(string(data), "'real'") {
		t.Fatalf("persistent connection missing from config file: %s", data)
	}
}
