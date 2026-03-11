package cmd

import (
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
		var symbols []string

		for _, symbol := range args {
			q, err := api.GetQuote(symbol)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error fetching %s: %v\n", symbol, err)
				continue
			}
			quotes = append(quotes, q)
			symbols = append(symbols, symbol)
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
