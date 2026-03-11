package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xgg",
	Short: "Xungungo CLI - Financial markets at your fingertips",
	Long:  "Xungungo CLI provides real-time stock quotes, historical data, and portfolio tracking from your terminal.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
