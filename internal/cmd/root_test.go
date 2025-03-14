package cmd

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that the command is properly configured
	if rootCmd.Use != "cg" {
		t.Errorf("Expected command use to be 'cg', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Command short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Command long description should not be empty")
	}

	if rootCmd.Run == nil {
		t.Error("Command run function should not be nil")
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Test that the upgrade flag is properly configured
	upgradeFlag := rootCmd.Flags().Lookup("upgrade")
	if upgradeFlag == nil {
		t.Error("Expected 'upgrade' flag to be defined")
		return
	}

	if upgradeFlag.Shorthand != "u" {
		t.Errorf("Expected shorthand for 'upgrade' flag to be 'u', got '%s'", upgradeFlag.Shorthand)
	}

	if upgradeFlag.Usage == "" {
		t.Error("Flag usage description should not be empty")
	}

	// Test that the version flag is properly configured
	versionFlag := rootCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Error("Expected 'version' flag to be defined")
		return
	}

	if versionFlag.Shorthand != "v" {
		t.Errorf("Expected shorthand for 'version' flag to be 'v', got '%s'", versionFlag.Shorthand)
	}

	if versionFlag.Usage == "" {
		t.Error("Flag usage description should not be empty")
	}
}

func TestRootCommandSubcommands(t *testing.T) {
	// Test that the upgrade subcommand is properly added
	upgradeFound := false
	versionFound := false

	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "upgrade" {
			upgradeFound = true
		}
		if cmd.Use == "version" {
			versionFound = true
		}
	}

	if !upgradeFound {
		t.Error("Expected 'upgrade' subcommand to be added to root command")
	}

	if !versionFound {
		t.Error("Expected 'version' subcommand to be added to root command")
	}
}
