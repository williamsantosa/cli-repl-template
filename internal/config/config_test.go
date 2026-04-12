package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func resetState(t *testing.T) {
	t.Helper()
	viper.Reset()
	C = Config{}
	t.Cleanup(func() {
		viper.Reset()
		C = Config{}
	})
}

func TestLoad_explicitFile(t *testing.T) {
	resetState(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.yaml")
	if err := os.WriteFile(path, []byte("name: mytool\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Load(path); err != nil {
		t.Fatal(err)
	}
	if C.Name != "mytool" {
		t.Fatalf("Name = %q, want mytool", C.Name)
	}
}

func TestLoad_defaultsWhenNoConfigInWd(t *testing.T) {
	resetState(t)
	tmp := t.TempDir()
	t.Chdir(tmp)
	if err := Load(""); err != nil {
		t.Fatal(err)
	}
	if C.Name != "cli-repl" {
		t.Fatalf("Name = %q, want default cli-repl", C.Name)
	}
}

func TestLoad_invalidYAML(t *testing.T) {
	resetState(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("art: [[[\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Load(path); err == nil {
		t.Fatal("Load: expected error for invalid YAML")
	}
}
