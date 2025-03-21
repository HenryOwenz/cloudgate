package constants

// Key constants for keyboard input
const (
	KeyQ          = "q"
	KeyCtrlC      = "ctrl+c"
	KeyEnter      = "enter"
	KeyCtrlEnter  = "ctrl+enter"
	KeyShiftEnter = "shift+enter"
	KeyF5         = "f5"
	KeyEsc        = "esc"
	KeyUp         = "up"
	KeyDown       = "down"
	KeyLeft       = "left"
	KeyRight      = "right"
	KeyAltUp      = "k"
	KeyAltDown    = "j"
	KeyAltBack    = "-"
	KeyTab        = "tab"

	// Vim-like navigation keys
	KeyGotoTop         = "g"
	KeyGotoBottom      = "G"
	KeyHome            = "home"
	KeyEnd             = "end"
	KeyHalfPageUp      = "ctrl+u"
	KeyHalfPageDown    = "ctrl+d"
	KeyAltHalfPageUp   = "u"
	KeyAltHalfPageDown = "d"
	KeyPageUp          = "pgup"
	KeyPageDown        = "pgdown"
	KeyAltPageUp       = "b"
	KeyAltPageDown     = "f"
	KeySpace           = " "

	// Pagination keys
	KeyNextPage          = "l"
	KeyPreviousPage      = "h"
	KeyArrowNextPage     = "right"
	KeyArrowPreviousPage = "left"

	// Search keys
	KeySearch    = "/"
	KeyBackspace = "backspace"
)

// Authentication method constants
const (
	// AWS authentication methods
	AWSProfileAuth = "profile"

	// Azure authentication methods (future)
	AzureCliAuth       = "cli"
	AzureConfigDirAuth = "config-dir"

	// GCP authentication methods (future)
	GCPServiceAccountAuth     = "service-account"
	GCPApplicationDefaultAuth = "adc"
)

// Configuration key constants
const (
	// AWS configuration keys
	AWSProfileKey = "profile"
	AWSRegionKey  = "region"

	// Azure configuration keys (future)
	AzureSubscriptionKey = "subscription"
	AzureLocationKey     = "location"
	AzureTenantKey       = "tenant"
	AzureConfigDirKey    = "config-dir"

	// GCP configuration keys (future)
	GCPProjectKey        = "project"
	GCPZoneKey           = "zone"
	GCPRegionKey         = "region"
	GCPServiceAccountKey = "service-account-path"
)
