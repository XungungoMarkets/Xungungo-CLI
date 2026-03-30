package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/output"
	"github.com/XungungoMarkets/xgg/internal/provider"
	"github.com/spf13/cobra"
)

var (
	screenerSector   string
	screenerCountry  string
	screenerIndustry string
)

var screenerCmd = &cobra.Command{
	Use:   "screener [symbol...]",
	Short: "Show raw NASDAQ screener data with all fields",
	Long:  "Fetch all NASDAQ-listed stocks and display the full screener table with every available field.",
	Example: `  xgg screener
  xgg screener AAPL MSFT
  xgg screener --sector Technology
  xgg screener --country USA
  xgg screener --industry "Software"
  xgg screener --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := provider.ServiceHandle()
		ctx := context.Background()

		rows, err := svc.GetScreenerRows(ctx)
		if err != nil {
			return fmt.Errorf("error fetching screener data: %w", err)
		}

		// Filter by symbol args (substring match on symbol or name).
		if len(args) > 0 {
			symbols := make(map[string]bool, len(args))
			for _, a := range args {
				symbols[strings.ToUpper(a)] = true
			}
			filtered := rows[:0]
			for _, r := range rows {
				if symbols[strings.ToUpper(r.Symbol)] {
					filtered = append(filtered, r)
				}
			}
			rows = filtered
		}

		if screenerSector != "" {
			filter := strings.ToLower(screenerSector)
			filtered := rows[:0]
			for _, r := range rows {
				if strings.Contains(strings.ToLower(r.Sector), filter) {
					filtered = append(filtered, r)
				}
			}
			rows = filtered
		}

		if screenerCountry != "" {
			filter := strings.ToLower(screenerCountry)
			filtered := rows[:0]
			for _, r := range rows {
				if strings.Contains(strings.ToLower(r.Country), filter) {
					filtered = append(filtered, r)
				}
			}
			rows = filtered
		}

		if screenerIndustry != "" {
			filter := strings.ToLower(screenerIndustry)
			filtered := rows[:0]
			for _, r := range rows {
				if strings.Contains(strings.ToLower(r.Industry), filter) {
					filtered = append(filtered, r)
				}
			}
			rows = filtered
		}

		if len(rows) == 0 {
			return fmt.Errorf("no stocks found matching the given filters")
		}

		if JSONOutput {
			jsonData, err := json.MarshalIndent(rows, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		output.PrintScreenerRows(rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(screenerCmd)
	screenerCmd.Flags().StringVar(&screenerSector, "sector", "", "Filter by sector (substring match)")
	screenerCmd.Flags().StringVar(&screenerCountry, "country", "", "Filter by country (substring match)")
	screenerCmd.Flags().StringVar(&screenerIndustry, "industry", "", "Filter by industry (substring match)")
}
