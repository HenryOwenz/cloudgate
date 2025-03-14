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
}

func TestRootCommandSubcommands(t *testing.T) {
	// Test that the upgrade subcommand is properly added
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "upgrade" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'upgrade' subcommand to be added to root command")
	}
}
