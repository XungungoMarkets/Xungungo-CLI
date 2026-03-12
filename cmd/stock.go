package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/XungungoMarkets/xgg/internal/api"
	"github.com/XungungoMarkets/xgg/internal/ui"
	"github.com/spf13/cobra"
)

var stockCmd = &cobra.Command{
	Use:   "stock [symbols...]",
	Short: "Get current stock quotes",
	Long:  "Fetch real-time stock quotes for one or more ticker symbols.",
	Example: `  xgg stock NVDA
  xgg stock NVDA AAPL TSLA
  xgg stock NVDA --json`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var quotes []*api.StockQuote
		service := api.ServiceHandle()

		for _, symbol := range args {
			q, meta, err := service.GetQuote(context.Background(), symbol)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error fetching %s: %v\n", symbol, err)
				continue
			}
			if meta.FallbackUsed && !JSONOutput {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: falling back to %s for %s (%v)\n", meta.ProviderUsed, symbol, meta.PrimaryErr)
			}
			quotes = append(quotes, q)
		}

		if len(quotes) == 0 {
			return fmt.Errorf("no quotes retrieved")
		}

		if JSONOutput {
			// Output in JSON format
			if len(quotes) == 1 {
				// Single quote
				jsonData, err := json.MarshalIndent(quotes[0], "", "  ")
				if err != nil {
					return fmt.Errorf("error marshaling JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			} else {
				// Multiple quotes
				jsonData, err := json.MarshalIndent(quotes, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshaling JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			}
		} else {
			// Output in human-readable format
			for _, q := range quotes {
				ui.PrintQuote(q)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
}
