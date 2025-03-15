package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Current is the current version of cloudgate
// This will be set during build time
var Current = "0.1.7"

// RepositoryURL is the base URL for the cloudgate repository
const RepositoryURL = "https://github.com/HenryOwenz/cloudgate"

// GitHubAPIURL is the URL for the GitHub API to fetch the latest release
const GitHubAPIURL = "https://api.github.com/repos/HenryOwenz/cloudgate/releases/latest"

// ReleaseInfo represents the GitHub API response for a release
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// GetLatest fetches the latest version from GitHub
func GetLatest() (string, error) {
	resp, err := http.Get(GitHubAPIURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var releaseInfo ReleaseInfo
	if err := json.Unmarshal(body, &releaseInfo); err != nil {
		return "", err
	}

	// Clean up the tag name (remove 'v' prefix if present)
	version := strings.TrimPrefix(releaseInfo.TagName, "v")

	return version, nil
}

// IsUpdateAvailable checks if a new version is available
func IsUpdateAvailable() (bool, string, error) {
	latestVersion, err := GetLatest()
	if err != nil {
		return false, "", err
	}

	// Compare versions (simple string comparison for now)
	// In a real implementation, you might want to use a proper version comparison library
	return latestVersion != Current, latestVersion, nil
}

// ColoredUpdateMessage returns a colored message about version status
func ColoredUpdateMessage() string {
	isNew, latestVersion, err := IsUpdateAvailable()
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
			yellow, cyan, Current, reset, cyan, latestVersion, reset, cyan, reset, RepositoryURL)
	}

	return ""
}
