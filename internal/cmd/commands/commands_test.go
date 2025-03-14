package commands

import (
	"runtime"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Test that the command is properly configured
	cmd := NewVersionCmd()

	if cmd.Use != "version" {
		t.Errorf("Expected command use to be 'version', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command short description should not be empty")
	}

	if cmd.Long == "" {
		t.Error("Command long description should not be empty")
	}

	if cmd.Run == nil {
		t.Error("Command run function should not be nil")
	}
}

func TestUpgradeCommand(t *testing.T) {
	// Test that the command is properly configured
	cmd := NewUpgradeCmd()

	if cmd.Use != "upgrade" {
		t.Errorf("Expected command use to be 'upgrade', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command short description should not be empty")
	}

	if cmd.Long == "" {
		t.Error("Command long description should not be empty")
	}

	if cmd.Run == nil {
		t.Error("Command run function should not be nil")
	}
}

func TestUpgradeCloudgateOSDetection(t *testing.T) {
	// Test that the OS detection works correctly
	os := runtime.GOOS

	switch os {
	case "windows", "darwin", "linux":
		t.Logf("OS %s is supported, skipping actual upgrade execution", os)
	default:
		err := upgradeCloudgate()
		if err == nil {
			t.Errorf("Expected error for unsupported OS: %s", os)
		}
	}
}
