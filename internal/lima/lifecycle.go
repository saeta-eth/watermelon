package lima

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VMStatus int

const (
	StatusNotFound VMStatus = iota
	StatusStopped
	StatusRunning
)

func (s VMStatus) String() string {
	switch s {
	case StatusRunning:
		return "Running"
	case StatusStopped:
		return "Stopped"
	default:
		return "Not found"
	}
}

// VMNameFromPath generates a consistent VM name from project path
func VMNameFromPath(projectPath string) string {
	base := filepath.Base(projectPath)
	base = strings.ToLower(base)
	base = strings.ReplaceAll(base, " ", "-")

	// Add short hash for uniqueness
	hash := sha256.Sum256([]byte(projectPath))
	shortHash := hex.EncodeToString(hash[:])[:8]

	return fmt.Sprintf("watermelon-%s-%s", base, shortHash)
}

// GetStatus returns the status of a VM
func GetStatus(vmName string) VMStatus {
	cmd := exec.Command("limactl", "list", "--format", "{{.Status}}", vmName)
	out, err := cmd.Output()
	if err != nil {
		return StatusNotFound
	}
	return parseStatus(strings.TrimSpace(string(out)))
}

func parseStatus(s string) VMStatus {
	switch s {
	case "Running":
		return StatusRunning
	case "Stopped":
		return StatusStopped
	default:
		return StatusNotFound
	}
}

// Start starts or creates a VM
func Start(vmName, configPath string) error {
	status := GetStatus(vmName)

	switch status {
	case StatusRunning:
		return nil // already running
	case StatusStopped:
		cmd := exec.Command("limactl", "start", vmName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	default:
		// Create new VM
		cmd := exec.Command("limactl", "start", "--name", vmName, configPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

// Stop stops a VM
func Stop(vmName string) error {
	cmd := exec.Command("limactl", "stop", vmName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Delete deletes a VM
func Delete(vmName string) error {
	cmd := exec.Command("limactl", "delete", "--force", vmName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Shell opens an interactive shell in the VM
func Shell(vmName string) error {
	cmd := exec.Command("limactl", "shell", vmName, "--", "bash", "-c", "cd /project && exec bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	// Ignore normal shell exit codes (0, 130=SIGINT, 143=SIGTERM)
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		if code == 0 || code == 130 || code == 143 {
			return nil
		}
	}
	return err
}

// Exec runs a command in the VM
func Exec(vmName string, args []string) error {
	// Wrap user command to run from /project directory
	shellCmd := fmt.Sprintf("cd /project && %s", strings.Join(args, " "))
	cmdArgs := []string{"shell", vmName, "--", "bash", "-c", shellCmd}
	cmd := exec.Command("limactl", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
