package test

import (
	"testing"
	"time"

	"github.com/XungungoMarkets/xgg/cmd"
	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/markcheno/go-talib"
)

func TestGetRSIData(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		bars      []market.Bar
		wantNil   bool
		wantValue float64
		wantType  string
	}{
		{
			name:    "Not enough data",
			symbol:  "AAPL",
			bars:    generateBars(10), // RSI needs 14
			wantNil: true,
		},
		{
			name:     "Valid RSI calculation",
			symbol:   "AAPL",
			bars:     generateBars(20),
			wantNil:  false,
			wantType: "rsi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cmd.GetRSIData(tt.symbol, tt.bars)
			if tt.wantNil && data != nil {
				t.Errorf("GetRSIData() should return nil, got %+v", data)
			}
			if !tt.wantNil && data == nil {
				t.Errorf("GetRSIData() should not return nil")
			}
			if data != nil {
				if data.Symbol != tt.symbol {
					t.Errorf("GetRSIData() symbol = %v, want %v", data.Symbol, tt.symbol)
				}
				if data.Indicator != "RSI(14)" {
					t.Errorf("GetRSIData() indicator = %v, want %v", data.Indicator, "RSI(14)")
				}
			}
		})
	}
}

func TestGetMACDData(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		bars    []market.Bar
		wantNil bool
	}{
		{
			name:    "Not enough data",
			symbol:  "AAPL",
			bars:    generateBars(20), // MACD needs 26
			wantNil: true,
		},
		{
			name:    "Valid MACD calculation",
			symbol:  "AAPL",
			bars:    generateBars(30),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cmd.GetMACDData(tt.symbol, tt.bars)
			if tt.wantNil && data != nil {
				t.Errorf("GetMACDData() should return nil, got %+v", data)
			}
			if !tt.wantNil && data == nil {
				t.Errorf("GetMACDData() should not return nil")
			}
			if data != nil {
				if data.Symbol != tt.symbol {
					t.Errorf("GetMACDData() symbol = %v, want %v", data.Symbol, tt.symbol)
				}
				if data.Indicator != "MACD(12,26,9)" {
					t.Errorf("GetMACDData() indicator = %v, want %v", data.Indicator, "MACD(12,26,9)")
				}
			}
		})
	}
}

func TestGetSMAData(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		bars    []market.Bar
		wantNil bool
	}{
		{
			name:    "Not enough data",
			symbol:  "AAPL",
			bars:    generateBars(40), // SMA needs 50
			wantNil: true,
		},
		{
			name:    "Valid SMA calculation",
			symbol:  "AAPL",
			bars:    generateBars(60),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cmd.GetSMAData(tt.symbol, tt.bars)
			if tt.wantNil && data != nil {
				t.Errorf("GetSMAData() should return nil, got %+v", data)
			}
			if !tt.wantNil && data == nil {
				t.Errorf("GetSMAData() should not return nil")
			}
			if data != nil {
				if data.Symbol != tt.symbol {
					t.Errorf("GetSMAData() symbol = %v, want %v", data.Symbol, tt.symbol)
				}
				if data.Indicator != "SMA(20,50)" {
					t.Errorf("GetSMAData() indicator = %v, want %v", data.Indicator, "SMA(20,50)")
				}
			}
		})
	}
}

func TestGetEMAData(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		bars    []market.Bar
		wantNil bool
	}{
		{
			name:    "Not enough data",
			symbol:  "AAPL",
			bars:    generateBars(20), // EMA needs 26
			wantNil: true,
		},
		{
			name:    "Valid EMA calculation",
			symbol:  "AAPL",
			bars:    generateBars(30),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cmd.GetEMAData(tt.symbol, tt.bars)
			if tt.wantNil && data != nil {
				t.Errorf("GetEMAData() should return nil, got %+v", data)
			}
			if !tt.wantNil && data == nil {
				t.Errorf("GetEMAData() should not return nil")
			}
			if data != nil {
				if data.Symbol != tt.symbol {
					t.Errorf("GetEMAData() symbol = %v, want %v", data.Symbol, tt.symbol)
				}
				if data.Indicator != "EMA(12,26)" {
					t.Errorf("GetEMAData() indicator = %v, want %v", data.Indicator, "EMA(12,26)")
				}
			}
		})
	}
}

func TestGetBollingerBandsData(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		bars    []market.Bar
		wantNil bool
	}{
		{
			name:    "Not enough data",
			symbol:  "AAPL",
			bars:    generateBars(15), // BB needs 20
			wantNil: true,
		},
		{
			name:    "Valid Bollinger Bands calculation",
			symbol:  "AAPL",
			bars:    generateBars(25),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cmd.GetBollingerBandsData(tt.symbol, tt.bars)
			if tt.wantNil && data != nil {
				t.Errorf("GetBollingerBandsData() should return nil, got %+v", data)
			}
			if !tt.wantNil && data == nil {
				t.Errorf("GetBollingerBandsData() should not return nil")
			}
			if data != nil {
				if data.Symbol != tt.symbol {
					t.Errorf("GetBollingerBandsData() symbol = %v, want %v", data.Symbol, tt.symbol)
				}
				if data.Indicator != "Bollinger Bands(20,2)" {
					t.Errorf("GetBollingerBandsData() indicator = %v, want %v", data.Indicator, "Bollinger Bands(20,2)")
				}
			}
		})
	}
}

func TestRSISignal(t *testing.T) {
	bars := generateBars(20)
	data := cmd.GetRSIData("TEST", bars)
	if data != nil {
		// Verify signal is one of the expected values
		if data.Signal != "overbought" && data.Signal != "oversold" && data.Signal != "neutral" {
			t.Errorf("RSI signal should be one of 'overbought', 'oversold', or 'neutral', got %s", data.Signal)
		}
		// Verify RSI value is in valid range (0-100)
		if data.Value < 0 || data.Value > 100 {
			t.Errorf("RSI value should be between 0 and 100, got %.2f", data.Value)
		}
	}
}

func TestEMATrend(t *testing.T) {
	bars := generateBars(30)
	data := cmd.GetEMAData("TEST", bars)
	if data != nil {
		// Verify trend is one of the expected values
		if data.Trend != "bullish" && data.Trend != "bearish" && data.Trend != "neutral" {
			t.Errorf("EMA trend should be one of 'bullish', 'bearish', or 'neutral', got %s", data.Trend)
		}
		// Verify EMA values are valid (positive)
		if data.EMA12 <= 0 || data.EMA26 <= 0 {
			t.Errorf("EMA values should be positive, got EMA12=%.2f, EMA26=%.2f", data.EMA12, data.EMA26)
		}
	}
}

// Helper function to generate test bar data
func generateBars(count int) []market.Bar {
	bars := make([]market.Bar, count)
	basePrice := 100.0
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < count; i++ {
		// Create realistic price movements
		change := float64(i%10-5) * 0.5 // -2.5 to +2.5
		price := basePrice + float64(i)*0.1 + change

		bars[i] = market.Bar{
			Date:   baseDate.AddDate(0, 0, i),
			Open:   price - 0.5,
			High:   price + 0.5,
			Low:    price - 1.0,
			Close:  price,
			Volume: 1000000 + i*10000,
		}
	}

	return bars
}

// Test with actual technical analysis library
func TestTechnicalIndicatorsIntegration(t *testing.T) {
	bars := generateBars(100)
	closes := make([]float64, len(bars))
	for i, bar := range bars {
		closes[i] = bar.Close
	}

	// Test RSI
	rsi := talib.Rsi(closes, 14)
	if len(rsi) == 0 {
		t.Error("RSI calculation returned empty slice")
	}

	// Test MACD
	macd, signal, hist := talib.Macd(closes, 12, 26, 9)
	if len(macd) == 0 || len(signal) == 0 || len(hist) == 0 {
		t.Error("MACD calculation returned empty slice(s)")
	}

	// Test SMA
	sma := talib.Sma(closes, 20)
	if len(sma) == 0 {
		t.Error("SMA calculation returned empty slice")
	}

	// Test EMA
	ema := talib.Ema(closes, 12)
	if len(ema) == 0 {
		t.Error("EMA calculation returned empty slice")
	}

	// Test Bollinger Bands
	upper, middle, lower := talib.BBands(closes, 20, 2.0, 2.0, 0)
	if len(upper) == 0 || len(middle) == 0 || len(lower) == 0 {
		t.Error("Bollinger Bands calculation returned empty slice(s)")
	}
}
