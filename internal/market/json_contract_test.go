package market

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestStockQuoteJSONContract(t *testing.T) {
	q := StockQuote{
		Symbol:        "NVDA",
		Name:          "NVIDIA Corporation",
		Price:         186.03,
		Change:        1.26,
		ChangePercent: 0.68,
		Volume:        138663908,
		MarketCap:     0,
	}

	raw, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("marshal stock quote: %v", err)
	}
	out := string(raw)
	for _, key := range []string{
		`"symbol"`,
		`"name"`,
		`"price"`,
		`"change"`,
		`"change_percent"`,
		`"volume"`,
		`"market_cap"`,
	} {
		if !strings.Contains(out, key) {
			t.Fatalf("missing key %s in %s", key, out)
		}
	}
}

func TestHistoryBarJSONContract(t *testing.T) {
	bar := Bar{
		Date:   time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC),
		Open:   100,
		High:   101,
		Low:    99,
		Close:  100.5,
		Volume: 1000,
	}
	raw, err := json.Marshal(bar)
	if err != nil {
		t.Fatalf("marshal bar: %v", err)
	}
	out := string(raw)
	for _, key := range []string{
		`"date"`,
		`"open"`,
		`"high"`,
		`"low"`,
		`"close"`,
		`"volume"`,
	} {
		if !strings.Contains(out, key) {
			t.Fatalf("missing key %s in %s", key, out)
		}
	}
}
