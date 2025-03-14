package cmd

import (
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Test that the command is properly configured
	if VersionCmd.Use != "version" {
		t.Errorf("Expected command use to be 'version', got '%s'", VersionCmd.Use)
	}

	if VersionCmd.Short == "" {
		t.Error("Command short description should not be empty")
	}

	if VersionCmd.Long == "" {
		t.Error("Command long description should not be empty")
	}

	if VersionCmd.Run == nil {
		t.Error("Command run function should not be nil")
	}
}

func TestVersionVariable(t *testing.T) {
	// Test that the Version variable is set
	if Version == "" {
		t.Error("Version variable should not be empty")
	}
}

func TestColoredVersionMessage(t *testing.T) {
	// We can't easily test the actual output without mocking the HTTP client,
	// but we can at least ensure the function doesn't panic
	_ = ColoredVersionMessage()
}
