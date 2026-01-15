package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCommandRequiresConfig(t *testing.T) {
	dir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	err = runRun()
	if err == nil {
		t.Error("expected error when no config exists")
	}
}

func TestRunCommandLoadsConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".watermelon.toml")

	config := `
[vm]
image = "ubuntu-22.04"

[network]
allow = []
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	// Just test that config loads without error
	// (actual VM operations would require Lima installed)
	cfg, err := loadProjectConfig(dir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.VM.Image != "ubuntu-22.04" {
		t.Errorf("expected ubuntu-22.04, got %s", cfg.VM.Image)
	}
}
