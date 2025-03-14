package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// UpgradeCmd represents the upgrade command
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cloudgate to the latest version",
	Long:  `Upgrade cloudgate to the latest version from GitHub releases.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if already on the latest version
		isNew, latestVersion, _, err := IsNewVersionAvailable()
		if err != nil {
			fmt.Println("Unable to check for the latest version. Proceeding with upgrade anyway...")
		} else if !isNew {
			fmt.Printf("You are already on the latest version (%s).\n", Version)
			return
		} else {
			fmt.Printf("Upgrading cloudgate from %s to %s...\n", Version, latestVersion)
		}

		err = upgradeCloudgate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error upgrading cloudgate: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Upgrade completed successfully!")
	},
}

// upgradeCloudgate runs the appropriate upgrade script based on the OS
func upgradeCloudgate() error {
	switch runtime.GOOS {
	case "windows":
		return upgradeWindows()
	case "darwin", "linux":
		return upgradeUnix()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// upgradeUnix runs the Unix (Linux/macOS) upgrade script
func upgradeUnix() error {
	// Using the exact command from README.md
	cmd := exec.Command("bash", "-c", "bash -c \"$(curl -fsSL https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.sh)\"")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// upgradeWindows runs the Windows upgrade script
func upgradeWindows() error {
	// Using the exact command from README.md
	powershellCmd := `Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.ps1'))`
	cmd := exec.Command("powershell", "-Command", powershellCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
