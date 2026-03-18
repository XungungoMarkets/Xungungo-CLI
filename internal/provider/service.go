package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/XungungoMarkets/xgg/internal/config"
	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/XungungoMarkets/xgg/internal/provider/nasdaq"
	"github.com/XungungoMarkets/xgg/internal/provider/yahoo"
)

// MarketDataProvider is the interface that all market data providers must implement.
type MarketDataProvider interface {
	Name() string
	GetQuote(ctx context.Context, symbol string) (*market.StockQuote, error)
	GetHistory(ctx context.Context, symbol, period string) ([]market.Bar, error)
}

// SearchProvider is the interface for providers that support symbol search.
type SearchProvider interface {
	Search(ctx context.Context, query string, limit int, includeMarketData bool) ([]market.SearchResult, error)
}

// SectorProvider is the interface for providers that support sector data.
type SectorProvider interface {
	GetSectors(ctx context.Context) ([]market.SectorSummary, error)
}

// Service orchestrates market data fetching across providers with fallback support.
type Service struct {
	mode     config.ProviderMode
	primary  MarketDataProvider
	fallback MarketDataProvider
}

func NewService(cfg config.RuntimeConfig) *Service {
	cfg = cfg.Normalized()
	return &Service{
		mode:     cfg.Provider,
		primary:  nasdaq.New(cfg),
		fallback: yahoo.New(),
	}
}

func (s *Service) GetQuote(ctx context.Context, symbol string) (*market.StockQuote, config.CallMeta, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, config.CallMeta{}, fmt.Errorf("symbol is required")
	}

	switch s.mode {
	case config.ProviderLegacy:
		q, err := s.fallback.GetQuote(ctx, symbol)
		return q, config.CallMeta{ProviderUsed: s.fallback.Name()}, err
	case config.ProviderNasdaq:
		q, err := s.primary.GetQuote(ctx, symbol)
		return q, config.CallMeta{ProviderUsed: s.primary.Name()}, err
	default:
		q, err := s.primary.GetQuote(ctx, symbol)
		if err == nil {
			return q, config.CallMeta{ProviderUsed: s.primary.Name()}, nil
		}
		if !config.IsRecoverable(err) {
			return nil, config.CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, err
		}
		fallbackQ, fallbackErr := s.fallback.GetQuote(ctx, symbol)
		if fallbackErr != nil {
			return nil, config.CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, fmt.Errorf("primary failed: %w; fallback failed: %v", err, fallbackErr)
		}
		return fallbackQ, config.CallMeta{ProviderUsed: s.fallback.Name(), FallbackUsed: true, PrimaryErr: err}, nil
	}
}

func (s *Service) GetHistory(ctx context.Context, symbol, period string) ([]market.Bar, config.CallMeta, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, config.CallMeta{}, fmt.Errorf("symbol is required")
	}

	switch s.mode {
	case config.ProviderLegacy:
		bars, err := s.fallback.GetHistory(ctx, symbol, period)
		return bars, config.CallMeta{ProviderUsed: s.fallback.Name()}, err
	case config.ProviderNasdaq:
		bars, err := s.primary.GetHistory(ctx, symbol, period)
		return bars, config.CallMeta{ProviderUsed: s.primary.Name()}, err
	default:
		bars, err := s.primary.GetHistory(ctx, symbol, period)
		if err == nil {
			return bars, config.CallMeta{ProviderUsed: s.primary.Name()}, nil
		}
		if !config.IsRecoverable(err) {
			return nil, config.CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, err
		}
		fallbackBars, fallbackErr := s.fallback.GetHistory(ctx, symbol, period)
		if fallbackErr != nil {
			return nil, config.CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, fmt.Errorf("primary failed: %w; fallback failed: %v", err, fallbackErr)
		}
		return fallbackBars, config.CallMeta{ProviderUsed: s.fallback.Name(), FallbackUsed: true, PrimaryErr: err}, nil
	}
}

func (s *Service) Search(ctx context.Context, query string, limit int, includeMarketData bool) ([]market.SearchResult, config.CallMeta, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, config.CallMeta{}, fmt.Errorf("query is required")
	}
	if limit <= 0 {
		limit = 10
	}

	if s.mode == config.ProviderLegacy {
		return nil, config.CallMeta{ProviderUsed: s.fallback.Name()}, fmt.Errorf("search is only available with nasdaq provider")
	}

	sp, ok := s.primary.(SearchProvider)
	if !ok {
		return nil, config.CallMeta{ProviderUsed: s.primary.Name()}, fmt.Errorf("provider %s does not support search", s.primary.Name())
	}

	results, err := sp.Search(ctx, query, limit, includeMarketData)
	if err != nil {
		return nil, config.CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, err
	}
	return results, config.CallMeta{ProviderUsed: s.primary.Name()}, nil
}

// --- Singleton runtime state ---

var (
	runtimeMu      sync.RWMutex
	currentConfig  = config.DefaultRuntimeConfig()
	currentService = NewService(currentConfig)
)

func ConfigureRuntime(cfg config.RuntimeConfig) {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	currentConfig = cfg.Normalized()
	currentService = NewService(currentConfig)
}

func RuntimeConfigSnapshot() config.RuntimeConfig {
	runtimeMu.RLock()
	defer runtimeMu.RUnlock()
	return currentConfig
}

func ServiceHandle() *Service {
	runtimeMu.RLock()
	defer runtimeMu.RUnlock()
	return currentService
}

func GetQuote(symbol string) (*market.StockQuote, error) {
	q, _, err := ServiceHandle().GetQuote(context.Background(), symbol)
	return q, err
}

func GetHistory(symbol, period string) ([]market.Bar, error) {
	bars, _, err := ServiceHandle().GetHistory(context.Background(), symbol, period)
	return bars, err
}

func (s *Service) GetSectors(ctx context.Context) ([]market.SectorSummary, error) {
	if s.mode == config.ProviderLegacy {
		return nil, fmt.Errorf("sectors command requires nasdaq provider")
	}
	sp, ok := s.primary.(SectorProvider)
	if !ok {
		return nil, fmt.Errorf("provider %s does not support sector data", s.primary.Name())
	}
	return sp.GetSectors(ctx)
}
