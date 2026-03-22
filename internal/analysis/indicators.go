package analysis

import (
	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/markcheno/go-talib"
)

func GetRSIData(symbol string, bars []market.Bar) *RSIOutput {
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

	return &RSIOutput{
		Symbol:    symbol,
		Indicator: "RSI(14)",
		Value:     lastRSI,
		Signal:    signal,
	}
}

func GetMACDData(symbol string, bars []market.Bar) *MACDOutput {
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

	return &MACDOutput{
		Symbol:     symbol,
		Indicator:  "MACD(12,26,9)",
		MACD:       lastMACD,
		Signal:     lastSignal,
		Histogram:  lastHist,
		SignalType: signalType,
	}
}

func GetSMAData(symbol string, bars []market.Bar) *SMAOutput {
	if len(bars) < 50 {
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
	lastSMA50 := sma50[len(sma50)-1]

	var lastSMA200 float64
	if len(bars) >= 200 {
		sma200 := talib.Sma(closes, 200)
		lastSMA200 = sma200[len(sma200)-1]
	}

	var trend string
	if currentPrice > lastSMA20 && lastSMA20 > lastSMA50 {
		trend = "uptrend"
	} else if currentPrice < lastSMA20 && lastSMA20 < lastSMA50 {
		trend = "downtrend"
	} else {
		trend = "consolidation"
	}

	return &SMAOutput{
		Symbol:       symbol,
		Indicator:    "SMA(20,50,200)",
		CurrentPrice: currentPrice,
		SMA20:        lastSMA20,
		SMA50:        lastSMA50,
		SMA200:       lastSMA200,
		Trend:        trend,
	}
}

func GetEMAData(symbol string, bars []market.Bar) *EMAOutput {
	if len(bars) < 26 {
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

	var lastEMA200 float64
	if len(bars) >= 200 {
		ema200 := talib.Ema(closes, 200)
		lastEMA200 = ema200[len(ema200)-1]
	}

	var trend string
	if lastEMA12 > lastEMA26 {
		trend = "bullish"
	} else if lastEMA12 < lastEMA26 {
		trend = "bearish"
	} else {
		trend = "neutral"
	}

	return &EMAOutput{
		Symbol:       symbol,
		Indicator:    "EMA(12,26,200)",
		CurrentPrice: currentPrice,
		EMA12:        lastEMA12,
		EMA26:        lastEMA26,
		EMA200:       lastEMA200,
		Trend:        trend,
	}
}

func GetBollingerBandsData(symbol string, bars []market.Bar) *BollingerBandsOutput {
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

	return &BollingerBandsOutput{
		Symbol:       symbol,
		Indicator:    "Bollinger Bands(20,2)",
		Upper:        lastUpper,
		Middle:       lastMiddle,
		Lower:        lastLower,
		CurrentPrice: currentPrice,
		Position:     position,
	}
}
