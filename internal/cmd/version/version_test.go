package version

import (
	"testing"
)

func TestVersionVariable(t *testing.T) {
	// Test that the Version variable is set
	if Current == "" {
		t.Error("Version variable should not be empty")
	}
}

func TestColoredUpdateMessage(t *testing.T) {
	// We can't easily test the actual output without mocking the HTTP client,
	// but we can at least ensure the function doesn't panic
	_ = ColoredUpdateMessage()
}
