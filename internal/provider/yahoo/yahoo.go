package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

		PreviousClose:           q.RegularMarketPreviousClose,
		MarketState:             string(q.MarketState),
		PreMarketPrice:          q.PreMarketPrice,
		PreMarketChange:         q.PreMarketChange,
		PreMarketChangePercent:  q.PreMarketChangePercent,
		PostMarketPrice:         q.PostMarketPrice,
		PostMarketChange:        q.PostMarketChange,
		PostMarketChangePercent: q.PostMarketChangePercent,
	}, nil
}

func (p *Provider) Search(_ context.Context, query string, limit int, _ bool) ([]market.SearchResult, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "query2.finance.yahoo.com",
		Path:   "/v1/finance/search",
	}
	q := u.Query()
	q.Set("q", strings.TrimSpace(query))
	q.Set("lang", "en-US")
	q.Set("region", "US")
	q.Set("quotesCount", fmt.Sprintf("%d", limit))
	q.Set("quotesQueryId", "tss_match_phrase_query")
	q.Set("enableNews", "false")
	q.Set("enableResearchReports", "false")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://finance.yahoo.com")
	req.Header.Set("Referer", "https://finance.yahoo.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Quotes []struct {
			Symbol    string `json:"symbol"`
			Longname  string `json:"longname"`
			Shortname string `json:"shortname"`
			QuoteType string `json:"quoteType"`
		} `json:"quotes"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	results := make([]market.SearchResult, 0, len(data.Quotes))
	for _, item := range data.Quotes {
		name := item.Longname
		if name == "" {
			name = item.Shortname
		}
		if item.Symbol == "" || name == "" {
			continue
		}
		results = append(results, market.SearchResult{
			Symbol: item.Symbol,
			Name:   name,
			Type:   item.QuoteType,
		})
	}
	return results, nil
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
