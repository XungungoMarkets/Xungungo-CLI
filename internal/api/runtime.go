package api

import (
	"context"
	"sync"
)

var (
	runtimeMu      sync.RWMutex
	currentConfig  = DefaultRuntimeConfig()
	currentService = NewService(currentConfig)
)

func ConfigureRuntime(cfg RuntimeConfig) {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	currentConfig = cfg.normalized()
	currentService = NewService(currentConfig)
}

func RuntimeConfigSnapshot() RuntimeConfig {
	runtimeMu.RLock()
	defer runtimeMu.RUnlock()
	return currentConfig
}

func ServiceHandle() *Service {
	runtimeMu.RLock()
	defer runtimeMu.RUnlock()
	return currentService
}

func GetQuote(symbol string) (*StockQuote, error) {
	q, _, err := ServiceHandle().GetQuote(context.Background(), symbol)
	return q, err
}

func GetHistory(symbol, period string) ([]Bar, error) {
	bars, _, err := ServiceHandle().GetHistory(context.Background(), symbol, period)
	return bars, err
}
