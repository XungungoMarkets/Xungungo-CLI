package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/api"
	"github.com/markcheno/go-talib"
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
		bars, err := api.GetHistory(args[0], technicalPeriod)
		if err != nil {
			return fmt.Errorf("error fetching history for %s: %w", args[0], err)
		}

		if len(bars) < 2 {
			return fmt.Errorf("not enough data points for technical analysis")
		}

		// Parse indicators (comma-separated)
		indicators := strings.Split(technicalIndicator, ",")

		// Map of indicator functions
		type indicatorFunc func(string, []api.Bar)
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
					results = append(results, getRSIData(args[0], bars))
					results = append(results, getMACDData(args[0], bars))
					results = append(results, getSMAData(args[0], bars))
					results = append(results, getEMAData(args[0], bars))
					results = append(results, getBollingerBandsData(args[0], bars))
					break
				}

				switch ind {
				case "rsi":
					results = append(results, getRSIData(args[0], bars))
				case "macd":
					results = append(results, getMACDData(args[0], bars))
				case "sma":
					results = append(results, getSMAData(args[0], bars))
				case "ema":
					results = append(results, getEMAData(args[0], bars))
				case "bb":
					results = append(results, getBollingerBandsData(args[0], bars))
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

// RSI functions
func calculateRSI(symbol string, bars []api.Bar) {
	data := getRSIData(symbol, bars)
	if data == nil {
		return
	}
	printRSIHumanReadable(data)
}

func getRSIData(symbol string, bars []api.Bar) *api.RSIOutput {
	if len(bars) < 14 {
		return nil
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	rsi := talib.Rsi(closes, 14)
	lastRSI := rsi[len(rsi)-1]

	var signal string
	if lastRSI > 70 {
		signal = "overbought"
	} else if lastRSI < 30 {
		signal = "oversold"
	} else {
		signal = "neutral"
	}

	return &api.RSIOutput{
		Symbol:    symbol,
		Indicator: "RSI(14)",
		Value:     lastRSI,
		Signal:    signal,
	}
}

func printRSIHumanReadable(data *api.RSIOutput) {
	fmt.Printf("📊 %s - RSI (14)\n", data.Symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Current RSI: %7.2f       │\n", data.Value)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	if data.Signal == "overbought" {
		fmt.Println("🔴 Overbought - Potential sell signal")
	} else if data.Signal == "oversold" {
		fmt.Println("🟢 Oversold - Potential buy signal")
	} else {
		fmt.Println("⚪ Neutral zone")
	}
}

// MACD functions
func calculateMACD(symbol string, bars []api.Bar) {
	data := getMACDData(symbol, bars)
	if data == nil {
		return
	}
	printMACDHumanReadable(data)
}

func getMACDData(symbol string, bars []api.Bar) *api.MACDOutput {
	if len(bars) < 26 {
		return nil
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	macd, signal, hist := talib.Macd(closes, 12, 26, 9)

	lastMACD := macd[len(macd)-1]
	lastSignal := signal[len(signal)-1]
	lastHist := hist[len(hist)-1]

	var signalType string
	if lastMACD > lastSignal && lastHist > 0 {
		signalType = "bullish"
	} else if lastMACD < lastSignal && lastHist < 0 {
		signalType = "bearish"
	} else {
		signalType = "neutral"
	}

	return &api.MACDOutput{
		Symbol:     symbol,
		Indicator:  "MACD(12,26,9)",
		MACD:       lastMACD,
		Signal:     lastSignal,
		Histogram:  lastHist,
		SignalType: signalType,
	}
}

func printMACDHumanReadable(data *api.MACDOutput) {
	fmt.Printf("📊 %s - MACD (12, 26, 9)\n", data.Symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ MACD Line:  %7.2f       │\n", data.MACD)
	fmt.Printf("│ Signal Line:%7.2f       │\n", data.Signal)
	fmt.Printf("│ Histogram:  %7.2f       │\n", data.Histogram)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	if data.SignalType == "bullish" {
		fmt.Println("🟢 Bullish - Buy signal")
	} else if data.SignalType == "bearish" {
		fmt.Println("🔴 Bearish - Sell signal")
	} else {
		fmt.Println("⚪ No clear signal")
	}
}

// SMA functions
func calculateSMA(symbol string, bars []api.Bar) {
	data := getSMAData(symbol, bars)
	if data == nil {
		return
	}
	printSMAHumanReadable(data)
}

func getSMAData(symbol string, bars []api.Bar) *api.SMAOutput {
	if len(bars) < 20 {
		return nil
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	sma20 := talib.Sma(closes, 20)
	sma50 := talib.Sma(closes, 50)

	currentPrice := bars[len(bars)-1].Close
	lastSMA20 := sma20[len(sma20)-1]
	var lastSMA50 float64
	if len(sma50) > 0 {
		lastSMA50 = sma50[len(sma50)-1]
	}

	var trend string
	if currentPrice > lastSMA20 && lastSMA20 > lastSMA50 {
		trend = "uptrend"
	} else if currentPrice < lastSMA20 && lastSMA20 < lastSMA50 {
		trend = "downtrend"
	} else {
		trend = "consolidation"
	}

	return &api.SMAOutput{
		Symbol:       symbol,
		Indicator:    "SMA(20,50)",
		CurrentPrice: currentPrice,
		SMA20:        lastSMA20,
		SMA50:        lastSMA50,
		Trend:        trend,
	}
}

func printSMAHumanReadable(data *api.SMAOutput) {
	fmt.Printf("📊 %s - Simple Moving Averages\n", data.Symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Current Price: $%7.2f    │\n", data.CurrentPrice)
	fmt.Printf("│ SMA 20:        $%7.2f    │\n", data.SMA20)
	if !math.IsNaN(data.SMA50) {
		fmt.Printf("│ SMA 50:        $%7.2f    │\n", data.SMA50)
	}
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	if data.Trend == "uptrend" {
		fmt.Println("🟢 Strong uptrend - Price above SMA20, SMA20 above SMA50")
	} else if data.Trend == "downtrend" {
		fmt.Println("🔴 Strong downtrend - Price below SMA20, SMA20 below SMA50")
	} else {
		fmt.Println("⚪ Consolidation or weak trend")
	}
}

// EMA functions
func calculateEMA(symbol string, bars []api.Bar) {
	data := getEMAData(symbol, bars)
	if data == nil {
		return
	}
	printEMAHumanReadable(data)
}

func getEMAData(symbol string, bars []api.Bar) *api.EMAOutput {
	if len(bars) < 12 {
		return nil
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	ema12 := talib.Ema(closes, 12)
	ema26 := talib.Ema(closes, 26)

	currentPrice := bars[len(bars)-1].Close
	lastEMA12 := ema12[len(ema12)-1]
	lastEMA26 := ema26[len(ema26)-1]

	var trend string
	if lastEMA12 > lastEMA26 {
		trend = "bullish"
	} else if lastEMA12 < lastEMA26 {
		trend = "bearish"
	} else {
		trend = "neutral"
	}

	return &api.EMAOutput{
		Symbol:       symbol,
		Indicator:    "EMA(12,26)",
		CurrentPrice: currentPrice,
		EMA12:        lastEMA12,
		EMA26:        lastEMA26,
		Trend:        trend,
	}
}

func printEMAHumanReadable(data *api.EMAOutput) {
	fmt.Printf("📊 %s - Exponential Moving Averages\n", data.Symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Current Price: $%7.2f    │\n", data.CurrentPrice)
	fmt.Printf("│ EMA 12:        $%7.2f    │\n", data.EMA12)
	fmt.Printf("│ EMA 26:        $%7.2f    │\n", data.EMA26)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	if data.Trend == "bullish" {
		fmt.Println("🟢 Bullish - EMA12 above EMA26")
	} else if data.Trend == "bearish" {
		fmt.Println("🔴 Bearish - EMA12 below EMA26")
	} else {
		fmt.Println("⚪ Neutral")
	}
}

// Bollinger Bands functions
func calculateBollingerBands(symbol string, bars []api.Bar) {
	data := getBollingerBandsData(symbol, bars)
	if data == nil {
		return
	}
	printBollingerBandsHumanReadable(data)
}

func getBollingerBandsData(symbol string, bars []api.Bar) *api.BollingerBandsOutput {
	if len(bars) < 20 {
		return nil
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	upper, middle, lower := talib.BBands(closes, 20, 2.0, 2.0, 0)

	currentPrice := bars[len(bars)-1].Close
	lastUpper := upper[len(upper)-1]
	lastMiddle := middle[len(middle)-1]
	lastLower := lower[len(lower)-1]

	var position string
	if currentPrice >= lastUpper {
		position = "above_upper"
	} else if currentPrice <= lastLower {
		position = "below_lower"
	} else {
		position = "within_bands"
	}

	return &api.BollingerBandsOutput{
		Symbol:       symbol,
		Indicator:    "Bollinger Bands(20,2)",
		Upper:        lastUpper,
		Middle:       lastMiddle,
		Lower:        lastLower,
		CurrentPrice: currentPrice,
		Position:     position,
	}
}

func printBollingerBandsHumanReadable(data *api.BollingerBandsOutput) {
	fmt.Printf("📊 %s - Bollinger Bands (20, 2)\n", data.Symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Upper Band:   $%7.2f    │\n", data.Upper)
	fmt.Printf("│ Middle (SMA):$%7.2f    │\n", data.Middle)
	fmt.Printf("│ Lower Band:   $%7.2f    │\n", data.Lower)
	fmt.Printf("│ Current Price:$%7.2f    │\n", data.CurrentPrice)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	if data.Position == "above_upper" {
		fmt.Println("🔴 Near or above upper band - Potential reversal or breakout")
	} else if data.Position == "below_lower" {
		fmt.Println("🟢 Near or below lower band - Potential reversal or breakdown")
	} else {
		fmt.Println("⚪ Price within bands - Normal trading range")
	}
}

func init() {
	technicalCmd.Flags().StringVarP(&technicalPeriod, "period", "p", "1m", "Time period: 5d, 1m, 3m, 6m, 1y, 5y")
	technicalCmd.Flags().StringVarP(&technicalIndicator, "indicator", "i", "rsi", "Indicator: rsi, macd, sma, ema, bb, all (comma-separated)")
	rootCmd.AddCommand(technicalCmd)
}
