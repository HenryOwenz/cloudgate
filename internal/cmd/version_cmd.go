package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of cloudgate",
	Long:  `Display the current version of cloudgate and check if a new version is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cloudgate version %s\n", Version)

		// Check for new version and display message if available
		isNew, latestVersion, _, err := IsNewVersionAvailable()
		if err == nil && isNew {
			fmt.Printf("A new version is available: %s\n", latestVersion)
			fmt.Println("Run 'cg --upgrade' to upgrade to the latest version.")
		}
	},
}
