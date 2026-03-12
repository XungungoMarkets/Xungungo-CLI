package api

import (
	"testing"

	nasdaqapi "github.com/XungungoMarkets/xgg-nasdaq-go/nasdaq"
)

func TestMapQuoteRowToStockQuote(t *testing.T) {
	row := &nasdaqapi.QuoteRow{
		Symbol:           "AAPL",
		Name:             "Apple Inc.",
		LastSalePrice:    "$123.45",
		NetChange:        "-1.55",
		PercentageChange: "-1.24%",
		Volume:           "1,234,567",
		MarketCap:        "3.1T",
	}

	got, err := mapQuoteRowToStockQuote(row, "")
	if err != nil {
		t.Fatalf("mapQuoteRowToStockQuote() error = %v", err)
	}
	if got.Symbol != "AAPL" || got.Name != "Apple Inc." {
		t.Fatalf("unexpected identity fields: %+v", got)
	}
	if got.Price != 123.45 || got.Change != -1.55 || got.ChangePercent != -1.24 {
		t.Fatalf("unexpected price/change fields: %+v", got)
	}
	if got.Volume != 1234567 || got.MarketCap != 3100000000000 {
		t.Fatalf("unexpected volume/marketcap fields: %+v", got)
	}
}

func TestMapQuoteRowToStockQuoteInvalidPrice(t *testing.T) {
	row := &nasdaqapi.QuoteRow{
		Symbol:        "AAPL",
		LastSalePrice: "N/A",
	}
	if _, err := mapQuoteRowToStockQuote(row, "AAPL"); err == nil {
		t.Fatal("expected parsing error")
	}
}

func TestPeriodToDateStrings(t *testing.T) {
	start, end := periodToDateStrings("1m")
	if len(start) != 10 || len(end) != 10 {
		t.Fatalf("unexpected date format start=%q end=%q", start, end)
	}
}

func TestMapSearchSuggestion(t *testing.T) {
	got := mapSearchSuggestion(nasdaqapi.SearchSuggestion{
		Symbol:      " NVDA ",
		Name:        " NVIDIA Corporation ",
		Type:        " STOCKS ",
		Description: " Semiconductors ",
	})

	if got.Symbol != "NVDA" || got.Name != "NVIDIA Corporation" || got.Type != "STOCKS" || got.Description != "Semiconductors" {
		t.Fatalf("unexpected mapped suggestion: %+v", got)
	}
}
