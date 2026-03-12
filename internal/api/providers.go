package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
)

type ProviderMode string

const (
	ProviderAuto   ProviderMode = "auto"
	ProviderNasdaq ProviderMode = "nasdaq"
	ProviderLegacy ProviderMode = "legacy"
)

type RuntimeConfig struct {
	Provider      ProviderMode
	RateLimit     int
	MaxRetries    int
	RetryDelaySec int
	TimeoutSec    int
	WatchlistType string
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		Provider:      ProviderAuto,
		RateLimit:     2,
		MaxRetries:    3,
		RetryDelaySec: 2,
		TimeoutSec:    30,
		WatchlistType: "Rv",
	}
}

func (c RuntimeConfig) normalized() RuntimeConfig {
	def := DefaultRuntimeConfig()
	if c.Provider == "" {
		c.Provider = def.Provider
	}
	if c.RateLimit <= 0 {
		c.RateLimit = def.RateLimit
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = def.MaxRetries
	}
	if c.RetryDelaySec <= 0 {
		c.RetryDelaySec = def.RetryDelaySec
	}
	if c.TimeoutSec <= 0 {
		c.TimeoutSec = def.TimeoutSec
	}
	if c.WatchlistType == "" {
		c.WatchlistType = def.WatchlistType
	}
	return c
}

type MarketDataProvider interface {
	Name() string
	GetQuote(ctx context.Context, symbol string) (*StockQuote, error)
	GetHistory(ctx context.Context, symbol, period string) ([]Bar, error)
}

type SearchProvider interface {
	Search(ctx context.Context, query string, limit int, includeMarketData bool) ([]SearchResult, error)
}

type CallMeta struct {
	ProviderUsed string
	FallbackUsed bool
	PrimaryErr   error
}

type Service struct {
	mode     ProviderMode
	primary  MarketDataProvider
	fallback MarketDataProvider
}

func NewService(cfg RuntimeConfig) *Service {
	cfg = cfg.normalized()
	return &Service{
		mode:     cfg.Provider,
		primary:  newNasdaqProvider(cfg),
		fallback: newLegacyProvider(),
	}
}

func (s *Service) GetQuote(ctx context.Context, symbol string) (*StockQuote, CallMeta, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, CallMeta{}, fmt.Errorf("symbol is required")
	}

	switch s.mode {
	case ProviderLegacy:
		q, err := s.fallback.GetQuote(ctx, symbol)
		return q, CallMeta{ProviderUsed: s.fallback.Name()}, err
	case ProviderNasdaq:
		q, err := s.primary.GetQuote(ctx, symbol)
		return q, CallMeta{ProviderUsed: s.primary.Name()}, err
	default:
		q, err := s.primary.GetQuote(ctx, symbol)
		if err == nil {
			return q, CallMeta{ProviderUsed: s.primary.Name()}, nil
		}
		if !IsRecoverable(err) {
			return nil, CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, err
		}
		fallbackQ, fallbackErr := s.fallback.GetQuote(ctx, symbol)
		if fallbackErr != nil {
			return nil, CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, fmt.Errorf("primary failed: %w; fallback failed: %v", err, fallbackErr)
		}
		return fallbackQ, CallMeta{ProviderUsed: s.fallback.Name(), FallbackUsed: true, PrimaryErr: err}, nil
	}
}

func (s *Service) GetHistory(ctx context.Context, symbol, period string) ([]Bar, CallMeta, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, CallMeta{}, fmt.Errorf("symbol is required")
	}

	switch s.mode {
	case ProviderLegacy:
		bars, err := s.fallback.GetHistory(ctx, symbol, period)
		return bars, CallMeta{ProviderUsed: s.fallback.Name()}, err
	case ProviderNasdaq:
		bars, err := s.primary.GetHistory(ctx, symbol, period)
		return bars, CallMeta{ProviderUsed: s.primary.Name()}, err
	default:
		bars, err := s.primary.GetHistory(ctx, symbol, period)
		if err == nil {
			return bars, CallMeta{ProviderUsed: s.primary.Name()}, nil
		}
		if !IsRecoverable(err) {
			return nil, CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, err
		}
		fallbackBars, fallbackErr := s.fallback.GetHistory(ctx, symbol, period)
		if fallbackErr != nil {
			return nil, CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, fmt.Errorf("primary failed: %w; fallback failed: %v", err, fallbackErr)
		}
		return fallbackBars, CallMeta{ProviderUsed: s.fallback.Name(), FallbackUsed: true, PrimaryErr: err}, nil
	}
}

func (s *Service) Search(ctx context.Context, query string, limit int, includeMarketData bool) ([]SearchResult, CallMeta, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, CallMeta{}, fmt.Errorf("query is required")
	}
	if limit <= 0 {
		limit = 10
	}

	if s.mode == ProviderLegacy {
		return nil, CallMeta{ProviderUsed: s.fallback.Name()}, fmt.Errorf("search is only available with nasdaq provider")
	}

	sp, ok := s.primary.(SearchProvider)
	if !ok {
		return nil, CallMeta{ProviderUsed: s.primary.Name()}, fmt.Errorf("provider %s does not support search", s.primary.Name())
	}

	results, err := sp.Search(ctx, query, limit, includeMarketData)
	if err != nil {
		return nil, CallMeta{ProviderUsed: s.primary.Name(), PrimaryErr: err}, err
	}
	return results, CallMeta{ProviderUsed: s.primary.Name()}, nil
}

type RecoverableError struct {
	Op  string
	Err error
}

func (e *RecoverableError) Error() string {
	return fmt.Sprintf("recoverable %s error: %v", e.Op, e.Err)
}

func (e *RecoverableError) Unwrap() error { return e.Err }

func recoverableError(op string, err error) error {
	if err == nil {
		return nil
	}
	return &RecoverableError{Op: op, Err: err}
}

func IsRecoverable(err error) bool {
	if err == nil {
		return false
	}
	var re *RecoverableError
	if errors.As(err, &re) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) && (ne.Timeout() || ne.Temporary()) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "parse") ||
		strings.Contains(msg, "invalid")
}
