package analysis

// Technical indicator structures for JSON output

// RSIOutput represents RSI indicator output
type RSIOutput struct {
	Symbol    string  `json:"symbol"`
	Indicator string  `json:"indicator"`
	Value     float64 `json:"value"`
	Signal    string  `json:"signal"` // "overbought", "oversold", "neutral"
}

// MACDOutput represents MACD indicator output
type MACDOutput struct {
	Symbol     string  `json:"symbol"`
	Indicator  string  `json:"indicator"`
	MACD       float64 `json:"macd"`
	Signal     float64 `json:"signal"`
	Histogram  float64 `json:"histogram"`
	SignalType string  `json:"signal_type"` // "bullish", "bearish", "neutral"
}

// SMAOutput represents SMA indicator output
type SMAOutput struct {
	Symbol       string  `json:"symbol"`
	Indicator    string  `json:"indicator"`
	CurrentPrice float64 `json:"current_price"`
	SMA20        float64 `json:"sma_20"`
	SMA50        float64 `json:"sma_50"`
	SMA200       float64 `json:"sma_200,omitempty"`
	Trend        string  `json:"trend"` // "uptrend", "downtrend", "consolidation"
}

// EMAOutput represents EMA indicator output
type EMAOutput struct {
	Symbol       string  `json:"symbol"`
	Indicator    string  `json:"indicator"`
	CurrentPrice float64 `json:"current_price"`
	EMA12        float64 `json:"ema_12"`
	EMA26        float64 `json:"ema_26"`
	EMA200       float64 `json:"ema_200,omitempty"`
	Trend        string  `json:"trend"` // "bullish", "bearish", "neutral"
}

// BollingerBandsOutput represents Bollinger Bands indicator output
type BollingerBandsOutput struct {
	Symbol       string  `json:"symbol"`
	Indicator    string  `json:"indicator"`
	Upper        float64 `json:"upper"`
	Middle       float64 `json:"middle"`
	Lower        float64 `json:"lower"`
	CurrentPrice float64 `json:"current_price"`
	Position     string  `json:"position"` // "above_upper", "below_lower", "within_bands"
}
