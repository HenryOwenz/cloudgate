package cmd

import (
	"fmt"
	"os"

	"github.com/HenryOwenz/cloudgate/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cg",
	Short: "A terminal-based application that unifies multi-cloud operations",
	Long: `cloudgate is a terminal-based application that unifies multi-cloud operations 
across AWS, Azure, and GCP.

Where your clouds converge.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if upgrade flag is set
		upgrade, _ := cmd.Flags().GetBool("upgrade")
		if upgrade {
			// Run the upgrade command
			UpgradeCmd.Run(cmd, args)
			return
		}

		// Check if version flag is set
		version, _ := cmd.Flags().GetBool("version")
		if version {
			// Run the version command
			VersionCmd.Run(cmd, args)
			return
		}

		// Default behavior - run the UI
		// Clear the screen using ANSI escape codes (works cross-platform)
		fmt.Print("\033[H\033[2J")

		// Create and run the program
		p := tea.NewProgram(ui.New())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// After the UI exits, check for new version and display message if available
		fmt.Print(ColoredVersionMessage())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add the upgrade flag to the root command
	rootCmd.Flags().BoolP("upgrade", "u", false, "Upgrade cloudgate to the latest version")

	// Add the version flag to the root command
	rootCmd.Flags().BoolP("version", "v", false, "Display the current version of cloudgate")

	// Add commands
	rootCmd.AddCommand(UpgradeCmd)
	rootCmd.AddCommand(VersionCmd)
}
