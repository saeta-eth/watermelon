//go:build e2e
// +build e2e

package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestE2EWorkflow(t *testing.T) {
	// Skip if Lima not installed
	if _, err := exec.LookPath("limactl"); err != nil {
		t.Skip("Lima not installed, skipping E2E test")
	}

	// Create temp project
	dir := t.TempDir()

	// Build watermelon
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(dir, "watermelon"), "./cmd/watermelon")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	wm := filepath.Join(dir, "watermelon")

	// Test init
	initCmd := exec.Command(wm, "init")
	initCmd.Dir = dir
	if out, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("init failed: %v\n%s", err, out)
	}

	configPath := filepath.Join(dir, ".watermelon.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config not created: %v", err)
	}

	// Test status (before VM exists)
	statusCmd := exec.Command(wm, "status")
	statusCmd.Dir = dir
	out, _ := statusCmd.CombinedOutput()
	t.Logf("status output: %s", out)
}
