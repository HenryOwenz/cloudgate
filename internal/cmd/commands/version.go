package commands

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cmd/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the current version of cloudgate",
		Long:  `Display the current version of cloudgate and check if a new version is available.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Print the current version
			fmt.Printf("cloudgate version %s\n", version.Current)

			// Check for new version and display message if available
			isNew, latestVersion, err := version.IsUpdateAvailable()
			if err == nil && isNew {
				// Use ANSI color codes for styling
				// Yellow text
				yellow := "\033[33m"
				// Cyan text
				cyan := "\033[36m"
				// Reset color
				reset := "\033[0m"

				fmt.Printf("\n%sA new release of cloudgate is available: %s%s%s → %s%s%s\nTo upgrade, run: %scg --upgrade%s\n%s\n",
					yellow, cyan, version.Current, reset, cyan, latestVersion, reset, cyan, reset, version.RepositoryURL)
			}
		},
	}

	return cmd
}
