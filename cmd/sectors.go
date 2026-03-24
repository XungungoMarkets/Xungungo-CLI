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

var byIndustry bool
var sectorByStock bool

var sectorsCmd = &cobra.Command{
	Use:   "sectors [sector name...]",
	Short: "Show % change by market sector",
	Long:  "Fetch all NASDAQ-listed stocks and show the average daily % change grouped by sector.",
	Example: `  xgg sectors
  xgg sectors --by-industry
  xgg sectors --by-stock
  xgg sectors --by-stock Energy
  xgg sectors --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := provider.ServiceHandle()
		ctx := context.Background()

		var filter string
		if len(args) > 0 {
			filter = strings.ToLower(strings.Join(args, " "))
		}

		if sectorByStock {
			data, err := svc.GetSectorStocks(ctx)
			if err != nil {
				return fmt.Errorf("error fetching sector stock data: %w", err)
			}
			if filter != "" {
				filtered := data[:0]
				for _, s := range data {
					if strings.Contains(strings.ToLower(s.Sector), filter) {
						filtered = append(filtered, s)
					}
				}
				data = filtered
			}
			if len(data) == 0 {
				return fmt.Errorf("no data found for sector %q", strings.Join(args, " "))
			}
			if JSONOutput {
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshaling JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			} else {
				output.PrintSectorStocks(data)
			}
			return nil
		}

		if byIndustry {
			industries, err := svc.GetIndustries(ctx)
			if err != nil {
				return fmt.Errorf("error fetching industry data: %w", err)
			}
			if len(industries) == 0 {
				return fmt.Errorf("no industry data available")
			}
			if JSONOutput {
				jsonData, err := json.MarshalIndent(industries, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshaling JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			} else {
				output.PrintIndustries(industries)
			}
			return nil
		}

		sectors, err := svc.GetSectors(ctx)
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
	sectorsCmd.Flags().BoolVar(&byIndustry, "by-industry", false, "Group by sector and industry")
	sectorsCmd.Flags().BoolVar(&sectorByStock, "by-stock", false, "Show individual stocks within each sector")
}
