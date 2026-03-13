package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/XungungoMarkets/xgg/internal/output"
	"github.com/XungungoMarkets/xgg/internal/provider"
	"github.com/spf13/cobra"
)

var historyPeriod string
var historyInterval string

var historyCmd = &cobra.Command{
	Use:   "history [symbol]",
	Short: "Get historical price data",
	Long:  "Fetch historical OHLCV data for a ticker symbol.",
	Example: `  xgg history NVDA
  xgg history NVDA --period 1y
  xgg history NVDA --period 1y --interval week
  xgg history NVDA --period 5y --interval month
  xgg history NVDA --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bars, meta, err := provider.ServiceHandle().GetHistory(context.Background(), args[0], historyPeriod)
		if err != nil {
			return fmt.Errorf("error fetching history for %s: %w", args[0], err)
		}
		if meta.FallbackUsed && !JSONOutput {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: falling back to %s for %s history (%v)\n", meta.ProviderUsed, args[0], meta.PrimaryErr)
		}

		bars = market.ApplyInterval(bars, historyInterval)

		if JSONOutput {
			jsonData, err := json.MarshalIndent(bars, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			output.PrintHistory(args[0], bars)
		}

		return nil
	},
}

func init() {
	historyCmd.Flags().StringVarP(&historyPeriod, "period", "p", "1m",
		"Time period: 1w, 2w, 5d, 1m, 2m, 3m, 6m, 9m, 1y, 2y, 3y, 5y, 10y, max")
	historyCmd.Flags().StringVar(&historyInterval, "interval", "day",
		"Bar interval: day, week, month")
	rootCmd.AddCommand(historyCmd)
}
