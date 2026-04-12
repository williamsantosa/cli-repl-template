package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/spf13/viper"

	"github.com/williamsantosa/cli-repl-template/internal/config"
)

func resetConfigGlobals(t *testing.T) {
	t.Helper()
	viper.Reset()
	config.C = config.Config{}
	t.Cleanup(func() {
		viper.Reset()
		config.C = config.Config{}
	})
}

func TestRenderFrames_missingEmbeddedAssetFallsBack(t *testing.T) {
	resetConfigGlobals(t)
	prev := EmbeddedAssets
	t.Cleanup(func() { EmbeddedAssets = prev })
	EmbeddedAssets = fstest.MapFS{}

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "c.yaml")
	yaml := "name: t\nart:\n  source: \"missing.png\"\n  width: 10\n"
	if err := os.WriteFile(cfgPath, []byte(yaml), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := config.Load(cfgPath); err != nil {
		t.Fatal(err)
	}

	frames := RenderFrames()
	if len(frames) != 1 {
		t.Fatalf("len(frames) = %d", len(frames))
	}
	if frames[0].Rendered == "" {
		t.Fatal("empty frame")
	}
	if !strings.Contains(frames[0].Rendered, "█") {
		t.Fatalf("expected built-in block art, got: %q", frames[0].Rendered)
	}
}

func TestRenderFrames_builtin(t *testing.T) {
	resetConfigGlobals(t)
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "c.yaml")
	if err := os.WriteFile(cfgPath, []byte("name: t\nart:\n  source: built-in\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := config.Load(cfgPath); err != nil {
		t.Fatal(err)
	}
	frames := RenderFrames()
	if len(frames) != 1 || !strings.Contains(frames[0].Rendered, "█") {
		t.Fatalf("%#v", frames)
	}
}
