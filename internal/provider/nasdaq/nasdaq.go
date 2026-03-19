package nasdaq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/XungungoMarkets/xgg/internal/config"
	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/XungungoMarkets/xgg/internal/parse"
	nasdaqapi "github.com/XungungoMarkets/xgg-nasdaq-go/nasdaq"
)

type Provider struct {
	cfg        config.RuntimeConfig
	client     *nasdaqapi.Client
	httpClient *http.Client

	rlMu       sync.Mutex
	nextReqUTC time.Time
}

func New(cfg config.RuntimeConfig) *Provider {
	cfg = cfg.Normalized()
	return &Provider{
		cfg: cfg,
		client: nasdaqapi.NewClient(
			nasdaqapi.WithRateLimit(cfg.RateLimit),
			nasdaqapi.WithMaxRetries(cfg.MaxRetries),
			nasdaqapi.WithRetryDelay(time.Duration(cfg.RetryDelaySec)*time.Second),
			nasdaqapi.WithWatchlistType(cfg.WatchlistType),
		),
		httpClient: &http.Client{Timeout: time.Duration(cfg.TimeoutSec) * time.Second},
	}
}

func (p *Provider) Name() string {
	return "nasdaq"
}

func (p *Provider) GetQuote(ctx context.Context, symbol string) (*market.StockQuote, error) {
	row, err := p.client.GetQuote(ctx, strings.ToUpper(symbol), nasdaqapi.SymbolTypeStock)
	if err == nil {
		quote, mapErr := mapQuoteRowToStockQuote(row, symbol)
		if mapErr == nil {
			return quote, nil
		}
		err = mapErr
	}

	// Fallback parsing path for changed watchlist response shapes.
	directRow, directErr := p.getQuoteViaWatchlist(ctx, symbol)
	if directErr == nil {
		quote, mapErr := mapQuoteRowToStockQuote(directRow, symbol)
		if mapErr == nil {
			return quote, nil
		}
		err = mapErr
	}

	quote, err := mapQuoteRowToStockQuote(row, symbol)
	if err != nil {
		return nil, config.NewRecoverableError("nasdaq_quote_parse", err)
	}
	return quote, nil
}

func (p *Provider) GetHistory(ctx context.Context, symbol string, period string) ([]market.Bar, error) {
	startDate, endDate := periodToDateStrings(period)
	u := url.URL{
		Scheme: "https",
		Host:   "api.nasdaq.com",
		Path:   fmt.Sprintf("/api/quote/%s/historical", strings.ToUpper(symbol)),
	}
	q := u.Query()
	q.Set("assetclass", "stocks")
	q.Set("fromdate", startDate)
	q.Set("todate", endDate)
	q.Set("limit", "9999")
	u.RawQuery = q.Encode()

	raw, err := p.getWithRetry(ctx, u.String())
	if err != nil {
		return nil, config.NewRecoverableError("nasdaq_history", err)
	}

	var resp historicalResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, config.NewRecoverableError("nasdaq_history_parse_json", err)
	}
	if resp.Status.RCode != 200 {
		return nil, config.NewRecoverableError("nasdaq_history_status", fmt.Errorf("status %d", resp.Status.RCode))
	}
	if len(resp.Data.TradesTable.Rows) == 0 {
		return nil, config.NewRecoverableError("nasdaq_history_empty", fmt.Errorf("empty rows"))
	}

	bars := make([]market.Bar, 0, len(resp.Data.TradesTable.Rows))
	for _, row := range resp.Data.TradesTable.Rows {
		d, err := time.Parse("01/02/2006", strings.TrimSpace(row.Date))
		if err != nil {
			return nil, config.NewRecoverableError("nasdaq_history_parse_date", err)
		}

		open, err := parse.ParseFloat(row.Open)
		if err != nil {
			return nil, config.NewRecoverableError("nasdaq_history_parse_open", err)
		}
		high, err := parse.ParseFloat(row.High)
		if err != nil {
			return nil, config.NewRecoverableError("nasdaq_history_parse_high", err)
		}
		low, err := parse.ParseFloat(row.Low)
		if err != nil {
			return nil, config.NewRecoverableError("nasdaq_history_parse_low", err)
		}
		closeVal, err := parse.ParseFloat(row.Close)
		if err != nil {
			return nil, config.NewRecoverableError("nasdaq_history_parse_close", err)
		}
		vol, err := parse.ParseInt(row.Volume)
		if err != nil {
			return nil, config.NewRecoverableError("nasdaq_history_parse_volume", err)
		}

		bars = append(bars, market.Bar{
			Date:   d,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closeVal,
			Volume: vol,
		})
	}

	sort.Slice(bars, func(i, j int) bool {
		return bars[i].Date.Before(bars[j].Date)
	})
	return bars, nil
}

func (p *Provider) GetSectors(ctx context.Context) ([]market.SectorSummary, error) {
	resp, err := p.client.GetScreenerStocks(ctx, true)
	if err != nil {
		return nil, err
	}

	type accum struct {
		sum   float64
		count int
	}
	sectors := make(map[string]*accum)

	for _, row := range resp.Rows {
		sector := strings.TrimSpace(row.Sector)
		if sector == "" {
			continue
		}
		pct, err := parse.ParseFloat(row.PercentageChange)
		if err != nil {
			continue
		}
		if _, ok := sectors[sector]; !ok {
			sectors[sector] = &accum{}
		}
		sectors[sector].sum += pct
		sectors[sector].count++
	}

	result := make([]market.SectorSummary, 0, len(sectors))
	for name, a := range sectors {
		result = append(result, market.SectorSummary{
			Sector:    name,
			AvgChange: a.sum / float64(a.count),
			Count:     a.count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].AvgChange > result[j].AvgChange
	})

	return result, nil
}

func (p *Provider) GetIndustries(ctx context.Context) ([]market.IndustrySummary, error) {
	resp, err := p.client.GetScreenerStocks(ctx, true)
	if err != nil {
		return nil, err
	}

	type accum struct {
		sum   float64
		count int
	}
	type key struct{ sector, industry string }
	groups := make(map[key]*accum)

	for _, row := range resp.Rows {
		sector := strings.TrimSpace(row.Sector)
		industry := strings.TrimSpace(row.Industry)
		if sector == "" || industry == "" {
			continue
		}
		pct, err := parse.ParseFloat(row.PercentageChange)
		if err != nil {
			continue
		}
		k := key{sector, industry}
		if _, ok := groups[k]; !ok {
			groups[k] = &accum{}
		}
		groups[k].sum += pct
		groups[k].count++
	}

	result := make([]market.IndustrySummary, 0, len(groups))
	for k, a := range groups {
		result = append(result, market.IndustrySummary{
			Sector:    k.sector,
			Industry:  k.industry,
			AvgChange: a.sum / float64(a.count),
			Count:     a.count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Sector != result[j].Sector {
			return result[i].Sector < result[j].Sector
		}
		return result[i].AvgChange > result[j].AvgChange
	})

	return result, nil
}

func (p *Provider) Search(ctx context.Context, query string, limit int, includeMarketData bool) ([]market.SearchResult, error) {
	resp, err := p.client.Search(ctx, strings.TrimSpace(query), limit, includeMarketData)
	if err == nil && resp != nil {
		results := make([]market.SearchResult, 0, len(resp.Data))
		for _, row := range resp.Data {
			results = append(results, mapSearchSuggestion(row))
		}
		return results, nil
	}

	// Fallback parsing path for current autosuggest shape (status.rCode + data[].metadata)
	results, directErr := p.searchViaAutosuggest(ctx, query, limit, includeMarketData)
	if directErr == nil {
		return results, nil
	}
	if err != nil {
		return nil, config.NewRecoverableError("nasdaq_search", err)
	}
	return nil, config.NewRecoverableError("nasdaq_search", directErr)
}

func mapQuoteRowToStockQuote(row *nasdaqapi.QuoteRow, fallbackSymbol string) (*market.StockQuote, error) {
	if row == nil {
		return nil, fmt.Errorf("nil quote row")
	}
	symbol := strings.ToUpper(strings.TrimSpace(row.Symbol))
	if symbol == "" {
		symbol = strings.ToUpper(strings.TrimSpace(fallbackSymbol))
	}
	if symbol == "" {
		return nil, fmt.Errorf("empty symbol")
	}

	price, err := parse.ParseFloat(row.LastSalePrice)
	if err != nil {
		return nil, fmt.Errorf("parse price: %w", err)
	}
	change, err := parse.ParseFloat(row.NetChange)
	if err != nil {
		return nil, fmt.Errorf("parse change: %w", err)
	}
	changePct, err := parse.ParseFloat(row.PercentageChange)
	if err != nil {
		return nil, fmt.Errorf("parse change percent: %w", err)
	}
	volume, err := parse.ParseInt(row.Volume)
	if err != nil {
		return nil, fmt.Errorf("parse volume: %w", err)
	}
	marketCap, err := parse.ParseInt64(row.MarketCap)
	if err != nil {
		return nil, fmt.Errorf("parse market cap: %w", err)
	}

	return &market.StockQuote{
		Symbol:        symbol,
		Name:          strings.TrimSpace(row.Name),
		Price:         price,
		Change:        change,
		ChangePercent: changePct,
		Volume:        volume,
		MarketCap:     marketCap,
	}, nil
}

func mapSearchSuggestion(s nasdaqapi.SearchSuggestion) market.SearchResult {
	return market.SearchResult{
		Symbol:      strings.TrimSpace(s.Symbol),
		Name:        strings.TrimSpace(s.Name),
		Type:        strings.TrimSpace(s.Type),
		Description: strings.TrimSpace(s.Description),
	}
}

func (p *Provider) searchViaAutosuggest(ctx context.Context, query string, limit int, includeMarketData bool) ([]market.SearchResult, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "www.nasdaq.com",
		Path:   "/ai-search/external/content-search-bff/v1/autosuggest",
	}
	q := u.Query()
	q.Set("query", strings.TrimSpace(query))
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("use_cache", "true")
	q.Set("include_market_data", fmt.Sprintf("%t", includeMarketData))
	u.RawQuery = q.Encode()

	raw, err := p.getWithRetry(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var response struct {
		Status struct {
			RCode int `json:"rCode"`
		} `json:"status"`
		Data []struct {
			SuggestedWord string `json:"suggestedWord"`
			Metadata      struct {
				Symbol      string `json:"symbol"`
				Name        string `json:"name"`
				Title       string `json:"title"`
				Asset       string `json:"asset"`
				Description string `json:"description"`
				DocType     string `json:"doc_type"`
				SubDocType  string `json:"sub_doc_type"`
			} `json:"metadata"`
		} `json:"data"`
	}

	if err := json.Unmarshal(raw, &response); err != nil {
		return nil, err
	}
	if response.Status.RCode != 200 {
		return nil, fmt.Errorf("status %d", response.Status.RCode)
	}

	results := make([]market.SearchResult, 0, len(response.Data))
	for _, row := range response.Data {
		name := strings.TrimSpace(row.Metadata.Name)
		if name == "" {
			name = strings.TrimSpace(row.Metadata.Title)
		}
		desc := strings.TrimSpace(row.Metadata.Description)
		if desc == "" {
			desc = strings.TrimSpace(row.SuggestedWord)
		}
		typ := strings.TrimSpace(row.Metadata.Asset)
		if typ == "" {
			typ = strings.TrimSpace(row.Metadata.DocType)
		}
		if row.Metadata.SubDocType != "" {
			typ = typ + ":" + strings.TrimSpace(row.Metadata.SubDocType)
		}

		results = append(results, market.SearchResult{
			Symbol:      strings.TrimSpace(row.Metadata.Symbol),
			Name:        name,
			Type:        typ,
			Description: desc,
		})
	}
	return results, nil
}

func (p *Provider) getWithRetry(ctx context.Context, endpoint string) ([]byte, error) {
	var lastErr error
	retries := p.cfg.MaxRetries
	if retries < 0 {
		retries = 0
	}

	for attempt := 0; attempt <= retries; attempt++ {
		if err := p.waitRateLimit(ctx); err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://www.nasdaq.com/")
		req.Header.Set("Origin", "https://www.nasdaq.com")

		resp, err := p.httpClient.Do(req)
		if err != nil {
			lastErr = err
		} else {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				lastErr = readErr
			} else if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				lastErr = fmt.Errorf("http status %d", resp.StatusCode)
			} else if resp.StatusCode >= 400 {
				return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, string(body))
			} else {
				return body, nil
			}
		}

		if attempt < retries {
			delay := time.Duration(p.cfg.RetryDelaySec) * time.Second
			select {
			case <-time.After(delay * time.Duration(attempt+1)):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	return nil, lastErr
}

func (p *Provider) waitRateLimit(ctx context.Context) error {
	rate := p.cfg.RateLimit
	if rate <= 0 {
		rate = 1
	}

	interval := time.Second / time.Duration(rate)
	p.rlMu.Lock()
	wait := time.Until(p.nextReqUTC)
	if wait < 0 {
		wait = 0
	}
	p.nextReqUTC = time.Now().Add(interval)
	p.rlMu.Unlock()

	if wait == 0 {
		return nil
	}
	select {
	case <-time.After(wait):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func periodToDateStrings(period string) (start string, end string) {
	now := time.Now()
	var from time.Time
	switch period {
	case "1w":
		from = now.AddDate(0, 0, -7)
	case "2w":
		from = now.AddDate(0, 0, -14)
	case "5d":
		from = now.AddDate(0, 0, -5)
	case "1m":
		from = now.AddDate(0, -1, 0)
	case "2m":
		from = now.AddDate(0, -2, 0)
	case "3m":
		from = now.AddDate(0, -3, 0)
	case "6m":
		from = now.AddDate(0, -6, 0)
	case "9m":
		from = now.AddDate(0, -9, 0)
	case "1y":
		from = now.AddDate(-1, 0, 0)
	case "2y":
		from = now.AddDate(-2, 0, 0)
	case "3y":
		from = now.AddDate(-3, 0, 0)
	case "5y":
		from = now.AddDate(-5, 0, 0)
	case "10y":
		from = now.AddDate(-10, 0, 0)
	case "max":
		from = now.AddDate(-20, 0, 0)
	default:
		from = now.AddDate(0, -1, 0)
	}
	return from.Format("2006-01-02"), now.Format("2006-01-02")
}

type historicalResponse struct {
	Status struct {
		RCode int `json:"rCode"`
	} `json:"status"`
	Data struct {
		TradesTable struct {
			Rows []struct {
				Date   string `json:"date"`
				Close  string `json:"close"`
				Volume string `json:"volume"`
				Open   string `json:"open"`
				High   string `json:"high"`
				Low    string `json:"low"`
			} `json:"rows"`
		} `json:"tradesTable"`
	} `json:"data"`
}

func (p *Provider) getQuoteViaWatchlist(ctx context.Context, symbol string) (*nasdaqapi.QuoteRow, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.nasdaq.com",
		Path:   "/api/quote/watchlist",
	}
	q := u.Query()
	q.Add("symbol", strings.ToLower(symbol)+"|stocks")
	if strings.TrimSpace(p.cfg.WatchlistType) != "" {
		q.Set("type", p.cfg.WatchlistType)
	}
	u.RawQuery = q.Encode()

	raw, err := p.getWithRetry(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var generic struct {
		Status struct {
			RCode int `json:"rCode"`
		} `json:"status"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, err
	}
	if generic.Status.RCode != 200 {
		return nil, fmt.Errorf("status %d", generic.Status.RCode)
	}

	// New shape: data.rows
	var rowsShape struct {
		Rows []struct {
			Symbol    string `json:"symbol"`
			Name      string `json:"name"`
			LastSale  string `json:"lastSale"`
			Change    string `json:"change"`
			PctChange string `json:"pctChange"`
			Volume    string `json:"volume"`
			MarketCap string `json:"marketCap"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(generic.Data, &rowsShape); err == nil && len(rowsShape.Rows) > 0 {
		r := rowsShape.Rows[0]
		return &nasdaqapi.QuoteRow{
			Symbol:           r.Symbol,
			Name:             r.Name,
			LastSalePrice:    r.LastSale,
			NetChange:        r.Change,
			PercentageChange: r.PctChange,
			Volume:           r.Volume,
			MarketCap:        r.MarketCap,
		}, nil
	}

	// Legacy shape: data is []QuoteRow
	var oldRows []nasdaqapi.QuoteRow
	if err := json.Unmarshal(generic.Data, &oldRows); err == nil && len(oldRows) > 0 {
		return &oldRows[0], nil
	}

	return nil, fmt.Errorf("unsupported watchlist response shape")
}
