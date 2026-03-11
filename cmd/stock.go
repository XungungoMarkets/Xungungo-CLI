package cmd

import (
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
  xgg stock NVDA AAPL TSLA`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, symbol := range args {
			q, err := api.GetQuote(symbol)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error fetching %s: %v\n", symbol, err)
				continue
			}
			ui.PrintQuote(q)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
}
