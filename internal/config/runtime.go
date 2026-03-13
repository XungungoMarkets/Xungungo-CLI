package config

import (
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

func (c RuntimeConfig) Normalized() RuntimeConfig {
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

type CallMeta struct {
	ProviderUsed string
	FallbackUsed bool
	PrimaryErr   error
}

type RecoverableError struct {
	Op  string
	Err error
}

func (e *RecoverableError) Error() string {
	return fmt.Sprintf("recoverable %s error: %v", e.Op, e.Err)
}

func (e *RecoverableError) Unwrap() error { return e.Err }

func NewRecoverableError(op string, err error) error {
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
