package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/XungungoMarkets/xgg/internal/output"
	"github.com/XungungoMarkets/xgg/internal/provider"
	"github.com/spf13/cobra"
)

var sectorsCmd = &cobra.Command{
	Use:   "sectors",
	Short: "Show % change by market sector",
	Long:  "Fetch all NASDAQ-listed stocks and show the average daily % change grouped by sector.",
	Example: `  xgg sectors
  xgg sectors --json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sectors, err := provider.ServiceHandle().GetSectors(context.Background())
		if err != nil {
			return fmt.Errorf("error fetching sector data: %w", err)
		}
		if len(sectors) == 0 {
			return fmt.Errorf("no sector data available")
		}

		if JSONOutput {
			jsonData, err := json.MarshalIndent(sectors, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			output.PrintSectors(sectors)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sectorsCmd)
}
