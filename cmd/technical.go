package cmd

import (
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
  xgg technical NVDA --indicator all --period 3m`,
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

		// Execute each indicator
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

		return nil
	},
}

func calculateRSI(symbol string, bars []api.Bar) {
	if len(bars) < 14 {
		fmt.Printf("⚠ Need at least 14 data points for RSI calculation (have %d)\n", len(bars))
		return
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	rsi := talib.Rsi(closes, 14)

	lastRSI := rsi[len(rsi)-1]
	fmt.Printf("📊 %s - RSI (14)\n", symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Current RSI: %7.2f       │\n", lastRSI)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	// RSI Analysis
	if lastRSI > 70 {
		fmt.Println("🔴 Overbought - Potential sell signal")
	} else if lastRSI < 30 {
		fmt.Println("🟢 Oversold - Potential buy signal")
	} else {
		fmt.Println("⚪ Neutral zone")
	}
}

func calculateMACD(symbol string, bars []api.Bar) {
	if len(bars) < 26 {
		fmt.Printf("⚠ Need at least 26 data points for MACD calculation (have %d)\n", len(bars))
		return
	}

	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	macd, signal, hist := talib.Macd(closes, 12, 26, 9)

	lastMACD := macd[len(macd)-1]
	lastSignal := signal[len(signal)-1]
	lastHist := hist[len(hist)-1]

	fmt.Printf("📊 %s - MACD (12, 26, 9)\n", symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ MACD Line:  %7.2f       │\n", lastMACD)
	fmt.Printf("│ Signal Line:%7.2f       │\n", lastSignal)
	fmt.Printf("│ Histogram:  %7.2f       │\n", lastHist)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	// MACD Analysis
	if lastMACD > lastSignal && lastHist > 0 {
		fmt.Println("🟢 Bullish - Buy signal")
	} else if lastMACD < lastSignal && lastHist < 0 {
		fmt.Println("🔴 Bearish - Sell signal")
	} else {
		fmt.Println("⚪ No clear signal")
	}
}

func calculateSMA(symbol string, bars []api.Bar) {
	if len(bars) < 20 {
		fmt.Printf("⚠ Need at least 20 data points for SMA calculation (have %d)\n", len(bars))
		return
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

	fmt.Printf("📊 %s - Simple Moving Averages\n", symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Current Price: $%7.2f    │\n", currentPrice)
	fmt.Printf("│ SMA 20:        $%7.2f    │\n", lastSMA20)
	if !math.IsNaN(lastSMA50) {
		fmt.Printf("│ SMA 50:        $%7.2f    │\n", lastSMA50)
	}
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	// SMA Analysis
	if currentPrice > lastSMA20 && lastSMA20 > lastSMA50 {
		fmt.Println("🟢 Strong uptrend - Price above SMA20, SMA20 above SMA50")
	} else if currentPrice < lastSMA20 && lastSMA20 < lastSMA50 {
		fmt.Println("🔴 Strong downtrend - Price below SMA20, SMA20 below SMA50")
	} else {
		fmt.Println("⚪ Consolidation or weak trend")
	}
}

func calculateEMA(symbol string, bars []api.Bar) {
	if len(bars) < 12 {
		fmt.Printf("⚠ Need at least 12 data points for EMA calculation (have %d)\n", len(bars))
		return
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

	fmt.Printf("📊 %s - Exponential Moving Averages\n", symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Current Price: $%7.2f    │\n", currentPrice)
	fmt.Printf("│ EMA 12:        $%7.2f    │\n", lastEMA12)
	fmt.Printf("│ EMA 26:        $%7.2f    │\n", lastEMA26)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	// EMA Analysis
	if lastEMA12 > lastEMA26 {
		fmt.Println("🟢 Bullish - EMA12 above EMA26")
	} else if lastEMA12 < lastEMA26 {
		fmt.Println("🔴 Bearish - EMA12 below EMA26")
	} else {
		fmt.Println("⚪ Neutral")
	}
}

func calculateBollingerBands(symbol string, bars []api.Bar) {
	if len(bars) < 20 {
		fmt.Printf("⚠ Need at least 20 data points for Bollinger Bands calculation (have %d)\n", len(bars))
		return
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

	fmt.Printf("📊 %s - Bollinger Bands (20, 2)\n", symbol)
	fmt.Printf("┌─────────────────────────────┐\n")
	fmt.Printf("│ Upper Band:   $%7.2f    │\n", lastUpper)
	fmt.Printf("│ Middle (SMA):$%7.2f    │\n", lastMiddle)
	fmt.Printf("│ Lower Band:   $%7.2f    │\n", lastLower)
	fmt.Printf("│ Current Price:$%7.2f    │\n", currentPrice)
	fmt.Printf("└─────────────────────────────┘\n")
	fmt.Println()

	// Bollinger Bands Analysis
	if currentPrice >= lastUpper {
		fmt.Println("🔴 Near or above upper band - Potential reversal or breakout")
	} else if currentPrice <= lastLower {
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
