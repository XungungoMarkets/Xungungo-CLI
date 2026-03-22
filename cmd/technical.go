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
var technicalInterval string

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
		period := technicalPeriod
		if minP := minPeriodForTechnical(technicalIndicator); minP != "" &&
			periodDaysTech(period) < periodDaysTech(minP) {
			period = minP
		}

		bars, err := provider.GetHistory(args[0], period)
		if err != nil {
			return fmt.Errorf("error fetching history for %s: %w", args[0], err)
		}

		bars = market.ApplyInterval(bars, technicalInterval)

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

// minPeriodForTechnical returns the minimum fetch period required to compute
// all requested indicators. Returns "" if the default period is sufficient.
func minPeriodForTechnical(indicator string) string {
	maxBars := 0
	for _, ind := range strings.Split(indicator, ",") {
		ind = strings.TrimSpace(ind)
		var need int
		switch ind {
		case "all":
			need = 200 // SMA(200) / EMA(200) is the largest
		case "rsi":
			need = 14
		case "macd":
			need = 35 // 26 slow EMA + 9 signal
		case "sma":
			need = 200
		case "ema":
			need = 200
		case "bb":
			need = 20
		}
		if need > maxBars {
			maxBars = need
		}
	}
	if maxBars == 0 {
		return ""
	}
	// 1.5× buffer to account for weekends and holidays
	calDays := int(float64(maxBars) * 1.5)
	for _, p := range []struct {
		s string
		d int
	}{
		{"1m", 30}, {"2m", 60}, {"3m", 90}, {"6m", 180},
		{"9m", 270}, {"1y", 365}, {"2y", 730},
	} {
		if p.d >= calDays {
			return p.s
		}
	}
	return "2y"
}

func periodDaysTech(p string) int {
	switch p {
	case "5d":
		return 5
	case "1w":
		return 7
	case "2w":
		return 14
	case "1m":
		return 30
	case "2m":
		return 60
	case "3m":
		return 90
	case "6m":
		return 180
	case "9m":
		return 270
	case "1y":
		return 365
	case "2y":
		return 730
	case "3y":
		return 1095
	case "5y":
		return 1825
	case "10y":
		return 3650
	case "max":
		return 7300
	default:
		return 30
	}
}

func init() {
	technicalCmd.Flags().StringVarP(&technicalPeriod, "period", "p", "1m",
		"Time period: 1w, 2w, 5d, 1m, 2m, 3m, 6m, 9m, 1y, 2y, 3y, 5y, 10y, max")
	technicalCmd.Flags().StringVarP(&technicalIndicator, "indicator", "i", "rsi",
		"Indicator: rsi, macd, sma, ema, bb, all (comma-separated)")
	technicalCmd.Flags().StringVar(&technicalInterval, "interval", "day",
		"Bar interval: day, week, month")
	rootCmd.AddCommand(technicalCmd)
}
