package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of xgg",
	Long:  "Display the current version of xgg CLI tool.",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan, color.Bold)
		white := color.New(color.FgWhite)

		cyan.Println("📈 Xungungo CLI")
		white.Printf("Version: %s\n", Version)
		white.Println("GitHub: https://github.com/XungungoMarkets/Xungungo-CLI")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
