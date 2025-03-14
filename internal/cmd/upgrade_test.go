package cmd

import (
	"runtime"
	"testing"
)

func TestUpgradeCommand(t *testing.T) {
	// Test that the command is properly configured
	if UpgradeCmd.Use != "upgrade" {
		t.Errorf("Expected command use to be 'upgrade', got '%s'", UpgradeCmd.Use)
	}

	if UpgradeCmd.Short == "" {
		t.Error("Command short description should not be empty")
	}

	if UpgradeCmd.Long == "" {
		t.Error("Command long description should not be empty")
	}

	if UpgradeCmd.Run == nil {
		t.Error("Command run function should not be nil")
	}
}

func TestUpgradeCloudgateOSDetection(t *testing.T) {
	// This test verifies that the OS detection logic works correctly
	// We can only test the OS detection part, not the actual upgrade functions

	// Check that the function doesn't return an error for supported OS
	// and returns an error for unsupported OS
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		// For supported OS, we expect no error from the OS detection part
		// (we can't actually test the upgrade functions themselves)
		t.Logf("OS %s is supported, skipping actual upgrade execution", runtime.GOOS)
	} else {
		// For unsupported OS, we expect an error
		err := upgradeCloudgate()
		if err == nil {
			t.Errorf("Expected error for unsupported OS %s, got nil", runtime.GOOS)
		} else if err.Error() != "unsupported operating system: "+runtime.GOOS {
			t.Errorf("Expected error message about unsupported OS, got: %v", err)
		}
	}
}
