package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Version is the current version of cloudgate
// This will be set during build time
var Version = "0.1.4"

// LatestReleaseInfo represents the GitHub API response for the latest release
type LatestReleaseInfo struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// GetLatestVersion fetches the latest version from GitHub
func GetLatestVersion() (string, string, error) {
	resp, err := http.Get("https://api.github.com/repos/HenryOwenz/cloudgate/releases/latest")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var releaseInfo LatestReleaseInfo
	if err := json.Unmarshal(body, &releaseInfo); err != nil {
		return "", "", err
	}

	// Clean up the tag name (remove 'v' prefix if present)
	version := strings.TrimPrefix(releaseInfo.TagName, "v")

	return version, releaseInfo.HTMLURL, nil
}

// IsNewVersionAvailable checks if a new version is available
func IsNewVersionAvailable() (bool, string, string, error) {
	latestVersion, url, err := GetLatestVersion()
	if err != nil {
		return false, "", "", err
	}

	// Compare versions (simple string comparison for now)
	// In a real implementation, you might want to use a proper version comparison library
	return latestVersion != Version, latestVersion, url, nil
}

// ColoredVersionMessage returns a colored message about version status
func ColoredVersionMessage() string {
	isNew, latestVersion, url, err := IsNewVersionAvailable()
	if err != nil {
		// Silently fail and don't show any message
		return ""
	}

	if isNew {
		// Use ANSI color codes for styling
		// Yellow text
		yellow := "\033[33m"
		// Cyan text
		cyan := "\033[36m"
		// Reset color
		reset := "\033[0m"

		return fmt.Sprintf("\n%sA new release of cloudgate is available: %s%s%s → %s%s%s\nTo upgrade, run: %scg --upgrade%s\n%s\n",
			yellow, cyan, Version, reset, cyan, latestVersion, reset, cyan, reset, url)
	}

	return ""
}
