package market

import (
	"time"
)

type StockQuote struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Volume        int     `json:"volume"`
	MarketCap     int64   `json:"market_cap"`
}

type Bar struct {
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int       `json:"volume"`
}

func PeriodToRange(period string) (int, int) {
	now := time.Now()
	end := int(now.Unix())

	var start time.Time
	switch period {
	case "5d":
		start = now.AddDate(0, 0, -5)
	case "1m":
		start = now.AddDate(0, -1, 0)
	case "3m":
		start = now.AddDate(0, -3, 0)
	case "6m":
		start = now.AddDate(0, -6, 0)
	case "1y":
		start = now.AddDate(-1, 0, 0)
	case "5y":
		start = now.AddDate(-5, 0, 0)
	default:
		start = now.AddDate(0, -1, 0)
	}

	return int(start.Unix()), end
}
