package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/analysis"
	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/XungungoMarkets/xgg/internal/output"
	"github.com/XungungoMarkets/xgg/internal/provider"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var technicalPeriod string
var technicalIndicator string

var technicalCmd = &cobra.Command{
	Use:   "technical [symbol]",
	Short: "Get technical analysis indicators",
	Long:  "Calculate and display technical analysis indicators for a ticker symbol.",
	Example: `  xgg technical NVDA
  xgg technical NVDA --indicator rsi --period 1m
  xgg technical NVDA --indicator rsi,macd --period 3m
  xgg technical NVDA --indicator all --period 3m
  xgg technical NVDA --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bars, err := provider.GetHistory(args[0], technicalPeriod)
		if err != nil {
			return fmt.Errorf("error fetching history for %s: %w", args[0], err)
		}

		if len(bars) < 2 {
			return fmt.Errorf("not enough data points for technical analysis")
		}

		// Parse indicators (comma-separated)
		indicators := strings.Split(technicalIndicator, ",")

		// Map of indicator functions
		type indicatorFunc func(string, []market.Bar)
		indicatorFuncs := map[string]indicatorFunc{
			"rsi":  calculateRSI,
			"macd": calculateMACD,
			"sma":  calculateSMA,
			"ema":  calculateEMA,
			"bb":   calculateBollingerBands,
		}

		if JSONOutput {
			// JSON output mode
			var results []interface{}

			for _, ind := range indicators {
				ind = strings.TrimSpace(ind)

				if ind == "all" {
					// If "all" is specified, get all indicators
					results = append(results, GetRSIData(args[0], bars))
					results = append(results, GetMACDData(args[0], bars))
					results = append(results, GetSMAData(args[0], bars))
					results = append(results, GetEMAData(args[0], bars))
					results = append(results, GetBollingerBandsData(args[0], bars))
					break
				}

				switch ind {
				case "rsi":
					results = append(results, GetRSIData(args[0], bars))
				case "macd":
					results = append(results, GetMACDData(args[0], bars))
				case "sma":
					results = append(results, GetSMAData(args[0], bars))
				case "ema":
					results = append(results, GetEMAData(args[0], bars))
				case "bb":
					results = append(results, GetBollingerBandsData(args[0], bars))
				default:
					return fmt.Errorf("unknown indicator: %s. Use: rsi, macd, sma, ema, bb, or all", ind)
				}
			}

			jsonData, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			// Human-readable output mode
			for i, ind := range indicators {
				ind = strings.TrimSpace(ind)

				if ind == "all" {
					// If "all" is specified, run all indicators
					indicatorFuncs["rsi"](args[0], bars)
					fmt.Println()
					indicatorFuncs["macd"](args[0], bars)
					fmt.Println()
					indicatorFuncs["sma"](args[0], bars)
					fmt.Println()
					indicatorFuncs["ema"](args[0], bars)
					fmt.Println()
					indicatorFuncs["bb"](args[0], bars)
					break
				}

				if fn, exists := indicatorFuncs[ind]; exists {
					if i > 0 {
						fmt.Println()
					}
					fn(args[0], bars)
				} else {
					return fmt.Errorf("unknown indicator: %s. Use: rsi, macd, sma, ema, bb, or all", ind)
				}
			}
		}

		return nil
	},
}

func calculateRSI(symbol string, bars []market.Bar) {
	data := GetRSIData(symbol, bars)
	if data == nil {
		yellow := color.New(color.FgYellow)
		yellow.Printf("⚠️  Not enough data for RSI (requires 14 bars, have %d)\n", len(bars))
		return
	}
	output.PrintRSI(data)
}

func calculateMACD(symbol string, bars []market.Bar) {
	data := GetMACDData(symbol, bars)
	if data == nil {
		yellow := color.New(color.FgYellow)
		yellow.Printf("⚠️  Not enough data for MACD (requires 26 bars, have %d)\n", len(bars))
		return
	}
	output.PrintMACD(data)
}

func calculateSMA(symbol string, bars []market.Bar) {
	data := GetSMAData(symbol, bars)
	if data == nil {
		yellow := color.New(color.FgYellow)
		yellow.Printf("⚠️  Not enough data for SMA (requires 50 bars, have %d)\n", len(bars))
		return
	}
	output.PrintSMA(data)
}

func calculateEMA(symbol string, bars []market.Bar) {
	data := GetEMAData(symbol, bars)
	if data == nil {
		yellow := color.New(color.FgYellow)
		yellow.Printf("⚠️  Not enough data for EMA (requires 26 bars, have %d)\n", len(bars))
		return
	}
	output.PrintEMA(data)
}

func calculateBollingerBands(symbol string, bars []market.Bar) {
	data := GetBollingerBandsData(symbol, bars)
	if data == nil {
		yellow := color.New(color.FgYellow)
		yellow.Printf("⚠️  Not enough data for Bollinger Bands (requires 20 bars, have %d)\n", len(bars))
		return
	}
	output.PrintBollingerBands(data)
}

// GetRSIData delegates to analysis package (kept for test compatibility)
func GetRSIData(symbol string, bars []market.Bar) *analysis.RSIOutput {
	return analysis.GetRSIData(symbol, bars)
}

// GetMACDData delegates to analysis package (kept for test compatibility)
func GetMACDData(symbol string, bars []market.Bar) *analysis.MACDOutput {
	return analysis.GetMACDData(symbol, bars)
}

// GetSMAData delegates to analysis package (kept for test compatibility)
func GetSMAData(symbol string, bars []market.Bar) *analysis.SMAOutput {
	return analysis.GetSMAData(symbol, bars)
}

// GetEMAData delegates to analysis package (kept for test compatibility)
func GetEMAData(symbol string, bars []market.Bar) *analysis.EMAOutput {
	return analysis.GetEMAData(symbol, bars)
}

// GetBollingerBandsData delegates to analysis package (kept for test compatibility)
func GetBollingerBandsData(symbol string, bars []market.Bar) *analysis.BollingerBandsOutput {
	return analysis.GetBollingerBandsData(symbol, bars)
}

func init() {
	technicalCmd.Flags().StringVarP(&technicalPeriod, "period", "p", "1m", "Time period: 5d, 1m, 3m, 6m, 1y, 5y")
	technicalCmd.Flags().StringVarP(&technicalIndicator, "indicator", "i", "rsi", "Indicator: rsi, macd, sma, ema, bb, all (comma-separated)")
	rootCmd.AddCommand(technicalCmd)
}
