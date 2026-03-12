package api

import (
	"context"
	"errors"
	"testing"
)

type mockProvider struct {
	name    string
	quote   *StockQuote
	history []Bar
	search  []SearchResult
	err     error
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) GetQuote(_ context.Context, _ string) (*StockQuote, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.quote, nil
}
func (m *mockProvider) GetHistory(_ context.Context, _, _ string) ([]Bar, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.history, nil
}
func (m *mockProvider) Search(_ context.Context, _ string, _ int, _ bool) ([]SearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.search, nil
}

func TestServiceFallbackOnRecoverableError(t *testing.T) {
	svc := &Service{
		mode:    ProviderAuto,
		primary: &mockProvider{name: "nasdaq", err: recoverableError("quote", errors.New("429"))},
		fallback: &mockProvider{name: "legacy", quote: &StockQuote{
			Symbol: "AAPL",
		}},
	}

	q, meta, err := svc.GetQuote(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}
	if q.Symbol != "AAPL" {
		t.Fatalf("unexpected quote: %+v", q)
	}
	if !meta.FallbackUsed || meta.ProviderUsed != "legacy" {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

func TestServiceNoFallbackOnNonRecoverableError(t *testing.T) {
	svc := &Service{
		mode:     ProviderAuto,
		primary:  &mockProvider{name: "nasdaq", err: errors.New("bad request")},
		fallback: &mockProvider{name: "legacy", quote: &StockQuote{Symbol: "AAPL"}},
	}

	_, _, err := svc.GetQuote(context.Background(), "AAPL")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceSearchNasdaqMode(t *testing.T) {
	svc := &Service{
		mode: ProviderNasdaq,
		primary: &mockProvider{
			name: "nasdaq",
			search: []SearchResult{
				{Symbol: "NVDA", Name: "NVIDIA", Type: "stocks"},
			},
		},
		fallback: &mockProvider{name: "legacy"},
	}

	results, meta, err := svc.Search(context.Background(), "NVDA", 10, false)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 1 || results[0].Symbol != "NVDA" {
		t.Fatalf("unexpected results: %+v", results)
	}
	if meta.ProviderUsed != "nasdaq" {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

func TestServiceSearchLegacyMode(t *testing.T) {
	svc := &Service{
		mode:     ProviderLegacy,
		primary:  &mockProvider{name: "nasdaq"},
		fallback: &mockProvider{name: "legacy"},
	}

	_, _, err := svc.Search(context.Background(), "NVDA", 10, false)
	if err == nil {
		t.Fatal("expected error for legacy mode search")
	}
}
