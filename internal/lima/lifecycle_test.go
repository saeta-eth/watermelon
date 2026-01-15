package lima

import (
	"testing"
)

func TestVMNameFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/Users/test/myproject", "watermelon-myproject"},
		{"/Users/test/my-project", "watermelon-my-project"},
		{"/Users/test/My Project", "watermelon-my-project"},
	}

	for _, tc := range tests {
		got := VMNameFromPath(tc.path)
		// Should start with watermelon-
		if got[:11] != "watermelon-" {
			t.Errorf("VMNameFromPath(%q) = %q, expected prefix 'watermelon-'", tc.path, got)
		}
	}
}

func TestVMStatus(t *testing.T) {
	// Test status parsing
	status := parseStatus("Running")
	if status != StatusRunning {
		t.Errorf("expected StatusRunning, got %v", status)
	}

	status = parseStatus("Stopped")
	if status != StatusStopped {
		t.Errorf("expected StatusStopped, got %v", status)
	}

	status = parseStatus("")
	if status != StatusNotFound {
		t.Errorf("expected StatusNotFound, got %v", status)
	}
}
