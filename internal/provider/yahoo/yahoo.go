package yahoo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/XungungoMarkets/xgg-finance-go/chart"
	"github.com/XungungoMarkets/xgg-finance-go/datetime"
	"github.com/XungungoMarkets/xgg-finance-go/quote"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "yahoo-finance"
}

func (p *Provider) GetQuote(_ context.Context, symbol string) (*market.StockQuote, error) {
	q, err := quote.Get(strings.ToUpper(symbol))
	if err != nil {
		return nil, fmt.Errorf("could not fetch quote for %s: %w", symbol, err)
	}

	return &market.StockQuote{
		Symbol:        q.Symbol,
		Name:          q.ShortName,
		Price:         q.RegularMarketPrice,
		Change:        q.RegularMarketChange,
		ChangePercent: q.RegularMarketChangePercent,
		Volume:        q.RegularMarketVolume,
		MarketCap:     0,
	}, nil
}

func (p *Provider) GetHistory(_ context.Context, symbol string, period string) ([]market.Bar, error) {
	interval := datetime.OneDay
	start, end := market.PeriodToRange(period)

	params := &chart.Params{
		Symbol:   strings.ToUpper(symbol),
		Interval: interval,
	}
	params.Start = datetime.FromUnix(start)
	params.End = datetime.FromUnix(end)

	var bars []market.Bar
	iter := chart.Get(params)
	for iter.Next() {
		b := iter.Bar()
		open, _ := b.Open.Float64()
		high, _ := b.High.Float64()
		low, _ := b.Low.Float64()
		close_, _ := b.Close.Float64()

		bars = append(bars, market.Bar{
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
