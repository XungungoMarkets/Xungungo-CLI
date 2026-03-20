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

var byStock bool

var countryCmd = &cobra.Command{
	Use:   "country [country name...]",
	Short: "Show % change by country",
	Long:  "Fetch all NASDAQ-listed stocks and show the average daily % change grouped by country. Optionally filter by country name.",
	Example: `  xgg country
  xgg country --by-stock
  xgg country --by-stock uruguay
  xgg country --by-stock hong kong
  xgg country --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := provider.ServiceHandle()
		ctx := context.Background()

		var filter string
		if len(args) > 0 {
			filter = strings.ToLower(strings.Join(args, " "))
		}

		if byStock {
			data, err := svc.GetCountryStocks(ctx)
			if err != nil {
				return fmt.Errorf("error fetching country stock data: %w", err)
			}
			if filter != "" {
				filtered := data[:0]
				for _, c := range data {
					if strings.Contains(strings.ToLower(c.Country), filter) {
						filtered = append(filtered, c)
					}
				}
				data = filtered
			}
			if len(data) == 0 {
				return fmt.Errorf("no data found for country %q", strings.Join(args, " "))
			}
			if JSONOutput {
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("error marshaling JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			} else {
				output.PrintCountryStocks(data)
			}
			return nil
		}

		countries, err := svc.GetCountries(ctx)
		if err != nil {
			return fmt.Errorf("error fetching country data: %w", err)
		}
		if filter != "" {
			filtered := countries[:0]
			for _, c := range countries {
				if strings.Contains(strings.ToLower(c.Country), filter) {
					filtered = append(filtered, c)
				}
			}
			countries = filtered
		}
		if len(countries) == 0 {
			return fmt.Errorf("no data found for country %q", strings.Join(args, " "))
		}

		if JSONOutput {
			jsonData, err := json.MarshalIndent(countries, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			output.PrintCountries(countries)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(countryCmd)
	countryCmd.Flags().BoolVar(&byStock, "by-stock", false, "Show individual stocks grouped by country")
}
