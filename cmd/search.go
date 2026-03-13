package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/XungungoMarkets/xgg/internal/output"
	"github.com/XungungoMarkets/xgg/internal/provider"
	"github.com/spf13/cobra"
)

var searchLimit int
var searchIncludeMarketData bool

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search symbols",
	Long:  "Search stocks, ETFs, indices, and other symbols using NASDAQ autosuggest.",
	Example: `  xgg search NVDA
  xgg search Apple --limit 20
  xgg search semiconductors --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, _, err := provider.ServiceHandle().Search(context.Background(), args[0], searchLimit, searchIncludeMarketData)
		if err != nil {
			return fmt.Errorf("error searching %q: %w", args[0], err)
		}

		if JSONOutput {
			jsonData, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		output.PrintSearchResults(args[0], results)
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 10, "Maximum number of search results")
	searchCmd.Flags().BoolVar(&searchIncludeMarketData, "market-data", false, "Include market data in search results when available")
	rootCmd.AddCommand(searchCmd)
}
