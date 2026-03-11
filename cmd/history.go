package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/XungungoMarkets/xgg/internal/api"
	"github.com/XungungoMarkets/xgg/internal/ui"
	"github.com/spf13/cobra"
)

var historyPeriod string

var historyCmd = &cobra.Command{
	Use:   "history [symbol]",
	Short: "Get historical price data",
	Long:  "Fetch historical OHLCV data for a ticker symbol.",
	Example: `  xgg history NVDA
  xgg history NVDA --period 1y
  xgg history NVDA --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bars, err := api.GetHistory(args[0], historyPeriod)
		if err != nil {
			return fmt.Errorf("error fetching history for %s: %w", args[0], err)
		}

		if JSONOutput {
			// Output in JSON format
			jsonData, err := json.MarshalIndent(bars, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			// Output in human-readable format
			ui.PrintHistory(args[0], bars)
		}

		return nil
	},
}

func init() {
	historyCmd.Flags().StringVarP(&historyPeriod, "period", "p", "1m", "Time period: 5d, 1m, 3m, 6m, 1y, 5y")
	rootCmd.AddCommand(historyCmd)
}
