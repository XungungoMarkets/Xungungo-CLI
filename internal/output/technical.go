package output

import (
	"fmt"
	"math"

	"github.com/XungungoMarkets/xgg/internal/analysis"
	"github.com/fatih/color"
)

func PrintRSI(data *analysis.RSIOutput) {
	if data == nil {
		color.New(color.FgYellow).Println("⚠️  Not enough data for RSI")
		return
	}
	fmt.Printf("📊 %s - RSI (14)\n", data.Symbol)
	fmt.Printf("┌────────────────────────────┐\n")
	fmt.Printf("│ Current RSI: %7.2f       │\n", data.Value)
	fmt.Printf("└────────────────────────────┘\n")
	fmt.Println()

	if data.Signal == "overbought" {
		fmt.Println("🔴 Overbought - Potential sell signal")
	} else if data.Signal == "oversold" {
		fmt.Println("🟢 Oversold - Potential buy signal")
	} else {
		fmt.Println("⚪ Neutral zone")
	}
}

func PrintMACD(data *analysis.MACDOutput) {
	if data == nil {
		color.New(color.FgYellow).Println("⚠️  Not enough data for MACD")
		return
	}
	fmt.Printf("📊 %s - MACD (12, 26, 9)\n", data.Symbol)
	fmt.Printf("┌────────────────────────────┐\n")
	fmt.Printf("│ MACD Line:  %7.2f        │\n", data.MACD)
	fmt.Printf("│ Signal Line:%7.2f        │\n", data.Signal)
	fmt.Printf("│ Histogram:  %7.2f        │\n", data.Histogram)
	fmt.Printf("└────────────────────────────┘\n")
	fmt.Println()

	if data.SignalType == "bullish" {
		fmt.Println("🟢 Bullish - Buy signal")
	} else if data.SignalType == "bearish" {
		fmt.Println("🔴 Bearish - Sell signal")
	} else {
		fmt.Println("⚪ No clear signal")
	}
}

func PrintSMA(data *analysis.SMAOutput) {
	if data == nil {
		color.New(color.FgYellow).Println("⚠️  Not enough data for SMA")
		return
	}
	fmt.Printf("📊 %s - Simple Moving Averages\n", data.Symbol)
	fmt.Printf("┌────────────────────────────┐\n")
	fmt.Printf("│ Current Price: $%7.2f    │\n", data.CurrentPrice)
	fmt.Printf("│ SMA 20:        $%7.2f    │\n", data.SMA20)
	if !math.IsNaN(data.SMA50) {
		fmt.Printf("│ SMA 50:        $%7.2f    │\n", data.SMA50)
	}
	if data.SMA200 != 0 {
		fmt.Printf("│ SMA 200:       $%7.2f    │\n", data.SMA200)
	}
	fmt.Printf("└────────────────────────────┘\n")
	fmt.Println()

	if data.Trend == "uptrend" {
		fmt.Println("🟢 Strong uptrend - Price above SMA20, SMA20 above SMA50")
	} else if data.Trend == "downtrend" {
		fmt.Println("🔴 Strong downtrend - Price below SMA20, SMA20 below SMA50")
	} else {
		fmt.Println("⚪ Consolidation or weak trend")
	}
}

func PrintEMA(data *analysis.EMAOutput) {
	if data == nil {
		color.New(color.FgYellow).Println("⚠️  Not enough data for EMA")
		return
	}
	fmt.Printf("📊 %s - Exponential Moving Averages\n", data.Symbol)
	fmt.Printf("┌────────────────────────────┐\n")
	fmt.Printf("│ Current Price: $%7.2f    │\n", data.CurrentPrice)
	fmt.Printf("│ EMA 12:        $%7.2f    │\n", data.EMA12)
	fmt.Printf("│ EMA 26:        $%7.2f    │\n", data.EMA26)
	if data.EMA200 != 0 {
		fmt.Printf("│ EMA 200:       $%7.2f    │\n", data.EMA200)
	}
	fmt.Printf("└────────────────────────────┘\n")
	fmt.Println()

	if data.Trend == "bullish" {
		fmt.Println("🟢 Bullish - EMA12 above EMA26")
	} else if data.Trend == "bearish" {
		fmt.Println("🔴 Bearish - EMA12 below EMA26")
	} else {
		fmt.Println("⚪ Neutral")
	}
}

func PrintBollingerBands(data *analysis.BollingerBandsOutput) {
	if data == nil {
		color.New(color.FgYellow).Println("⚠️  Not enough data for Bollinger Bands")
		return
	}
	fmt.Printf("📊 %s - Bollinger Bands (20, 2)\n", data.Symbol)
	fmt.Printf("┌────────────────────────────┐\n")
	fmt.Printf("│ Upper Band:   $%7.2f     │\n", data.Upper)
	fmt.Printf("│ Middle (SMA): $%7.2f     │\n", data.Middle)
	fmt.Printf("│ Lower Band:   $%7.2f     │\n", data.Lower)
	fmt.Printf("│ Current Price:$%7.2f     │\n", data.CurrentPrice)
	fmt.Printf("└────────────────────────────┘\n")
	fmt.Println()

	if data.Position == "above_upper" {
		fmt.Println("🔴 Near or above upper band - Potential reversal or breakout")
	} else if data.Position == "below_lower" {
		fmt.Println("🟢 Near or below lower band - Potential reversal or breakdown")
	} else {
		fmt.Println("⚪ Price within bands - Normal trading range")
	}
}
