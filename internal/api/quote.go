package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/XungungoMarkets/xgg-finance-go/chart"
	"github.com/XungungoMarkets/xgg-finance-go/datetime"
	"github.com/XungungoMarkets/xgg-finance-go/quote"
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

func GetQuote(symbol string) (*StockQuote, error) {
	q, err := quote.Get(strings.ToUpper(symbol))
	if err != nil {
		return nil, fmt.Errorf("could not fetch quote for %s: %w", symbol, err)
	}

	return &StockQuote{
		Symbol:        q.Symbol,
		Name:          q.ShortName,
		Price:         q.RegularMarketPrice,
		Change:        q.RegularMarketChange,
		ChangePercent: q.RegularMarketChangePercent,
		Volume:        q.RegularMarketVolume,
		MarketCap:     0, // TODO: fetch from equity endpoint
	}, nil
}

func GetHistory(symbol string, period string) ([]Bar, error) {
	interval := datetime.OneDay

	start, end := periodToRange(period)

	params := &chart.Params{
		Symbol:   strings.ToUpper(symbol),
		Interval: interval,
	}
	params.Start = datetime.FromUnix(start)
	params.End = datetime.FromUnix(end)

	var bars []Bar
	iter := chart.Get(params)
	for iter.Next() {
		b := iter.Bar()
		open, _ := b.Open.Float64()
		high, _ := b.High.Float64()
		low, _ := b.Low.Float64()
		close_, _ := b.Close.Float64()

		bars = append(bars, Bar{
			Date:   time.Unix(int64(b.Timestamp), 0),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close_,
			Volume: b.Volume,
		})
	}
	if iter.Err() != nil {
		return nil, fmt.Errorf("could not fetch history for %s: %w", symbol, iter.Err())
	}

	return bars, nil
}

func periodToRange(period string) (int, int) {
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
