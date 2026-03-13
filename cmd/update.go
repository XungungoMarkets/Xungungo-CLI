package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/fatih/color"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const (
	repo = "XungungoMarkets/Xungungo-CLI"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update xgg to the latest version",
	Long:  "Check for updates and automatically update xgg to the latest version from GitHub releases.",
	Run: func(cmd *cobra.Command, args []string) {
		doUpdate()
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check-update",
	Short: "Check if an update is available",
	Long:  "Check if a newer version of xgg is available without installing it.",
	Run: func(cmd *cobra.Command, args []string) {
		checkUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(checkUpdateCmd)
}

func checkUpdate() {
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		log.Println("Error detecting version:", err)
		return
	}

	currentVersion, err := semverParse(Version)
	if err != nil {
		log.Println("Error parsing current version:", err)
		return
	}

	if !found {
		green := color.New(color.FgGreen)
		green.Println("✓ No releases found on GitHub")
		return
	}

	if latest.Version.LTE(currentVersion) {
		green := color.New(color.FgGreen)
		green.Printf("✓ You are using the latest version: %s\n", Version)
		return
	}

	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)

	yellow.Println("⚠ A new version is available!")
	cyan.Printf("  Current: %s\n", Version)
	cyan.Printf("  Latest:  %s\n", latest.Version.String())
	fmt.Printf("  Release:  %s\n\n", latest.AssetURL)

	white := color.New(color.FgWhite)
	white.Println("Run 'xgg update' to update to the latest version.")
}

func doUpdate() {
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		log.Println("Error detecting version:", err)
		os.Exit(1)
	}

	currentVersion, err := semverParse(Version)
	if err != nil {
		log.Println("Error parsing current version:", err)
		os.Exit(1)
	}

	if !found {
		log.Println("No releases found on GitHub")
		return
	}

	if latest.Version.LTE(currentVersion) {
		green := color.New(color.FgGreen)
		green.Println("✓ You are already using the latest version:", Version)
		return
	}

	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)

	yellow.Println("⚠ New version available!")
	cyan.Printf("  Current: %s\n", Version)
	cyan.Printf("  Latest:  %s\n", latest.Version.String())
	fmt.Printf("  Release:  %s\n\n", latest.AssetURL)

	if !confirm("Do you want to update?") {
		white := color.New(color.FgWhite)
		white.Println("Update cancelled.")
		return
	}

	cyan.Println("\n⏳ Downloading update...")

	execPath, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		os.Exit(1)
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, execPath); err != nil {
		log.Println("Error updating:", err)
		os.Exit(1)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Successfully updated to version", latest.Version.String())
	white := color.New(color.FgWhite)
	white.Println("Please restart xgg to use the new version.")
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/N]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		response = strings.TrimSpace(strings.ToLower(response))
		return response == "y" || response == "yes"
	}
}

// semverParse is a wrapper around semver.Make that handles "dev" version and "v" prefix
func semverParse(version string) (semver.Version, error) {
	if version == "dev" {
		return semver.Version{Major: 0, Minor: 0, Patch: 0}, nil
	}
	// Remove "v" prefix if present (e.g., "v0.2.2" -> "0.2.2")
	cleanVersion := strings.TrimPrefix(version, "v")
	return semver.Make(cleanVersion)
}
