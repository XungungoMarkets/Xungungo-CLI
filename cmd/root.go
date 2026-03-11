package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time using -ldflags
var Version = "dev"

// JSONOutput is a global flag to control JSON output format
var JSONOutput bool

var rootCmd = &cobra.Command{
	Use:   "xgg",
	Short: "Xungungo CLI - Financial markets at your fingertips",
	Long:  "Xungungo CLI provides real-time stock quotes, historical data, and portfolio tracking from your terminal.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set JSON output flag based on parent command
		if cmd.Parent() != nil {
			if JSONOutput, _ = cmd.Parent().Flags().GetBool("json"); !JSONOutput {
				JSONOutput, _ = cmd.Flags().GetBool("json")
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&JSONOutput, "json", false, "Output in JSON format")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
