package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/XungungoMarkets/xgg/internal/api"
	"github.com/spf13/cobra"
)

func newTestRootCmd() *cobra.Command {
	c := &cobra.Command{Use: "xgg"}
	c.PersistentFlags().StringVar(&providerMode, "provider", "", "")
	c.PersistentFlags().IntVar(&rateLimit, "rate-limit", 0, "")
	c.PersistentFlags().IntVar(&maxRetries, "max-retries", -1, "")
	c.PersistentFlags().DurationVar(&retryDelay, "retry-delay", 0, "")
	c.PersistentFlags().DurationVar(&requestTimeout, "timeout", 0, "")
	c.PersistentFlags().StringVar(&watchlistType, "watchlist-type", "", "")
	return c
}

func TestResolveRuntimeConfigEnv(t *testing.T) {
	t.Setenv("XGG_PROVIDER", "nasdaq")
	t.Setenv("XGG_RATE_LIMIT", "1")
	t.Setenv("XGG_MAX_RETRIES", "5")
	t.Setenv("XGG_RETRY_DELAY", "3s")
	t.Setenv("XGG_TIMEOUT", "40s")
	t.Setenv("XGG_WATCHLIST_TYPE", "Rv")

	cmd := newTestRootCmd()
	cfg := resolveRuntimeConfig(cmd)
	if cfg.Provider != api.ProviderNasdaq || cfg.RateLimit != 1 || cfg.MaxRetries != 5 {
		t.Fatalf("unexpected config from env: %+v", cfg)
	}
	if cfg.RetryDelaySec != 3 || cfg.TimeoutSec != 40 || cfg.WatchlistType != "Rv" {
		t.Fatalf("unexpected duration/watchlist config from env: %+v", cfg)
	}
}

func TestResolveRuntimeConfigFlagsOverrideEnv(t *testing.T) {
	t.Setenv("XGG_PROVIDER", "legacy")
	t.Setenv("XGG_RATE_LIMIT", "1")

	cmd := newTestRootCmd()
	if err := cmd.PersistentFlags().Set("provider", "nasdaq"); err != nil {
		t.Fatalf("set provider flag: %v", err)
	}
	if err := cmd.PersistentFlags().Set("rate-limit", "7"); err != nil {
		t.Fatalf("set rate-limit flag: %v", err)
	}
	if err := cmd.PersistentFlags().Set("retry-delay", "9s"); err != nil {
		t.Fatalf("set retry-delay flag: %v", err)
	}
	if err := cmd.PersistentFlags().Set("timeout", "22s"); err != nil {
		t.Fatalf("set timeout flag: %v", err)
	}

	cfg := resolveRuntimeConfig(cmd)
	if cfg.Provider != api.ProviderNasdaq || cfg.RateLimit != 7 {
		t.Fatalf("flags did not override env: %+v", cfg)
	}
	if cfg.RetryDelaySec != int((9*time.Second)/time.Second) || cfg.TimeoutSec != int((22*time.Second)/time.Second) {
		t.Fatalf("duration flags not applied: %+v", cfg)
	}
}

func TestResolveRuntimeConfigDefaults(t *testing.T) {
	for _, key := range []string{
		"XGG_PROVIDER",
		"XGG_RATE_LIMIT",
		"XGG_MAX_RETRIES",
		"XGG_RETRY_DELAY",
		"XGG_TIMEOUT",
		"XGG_WATCHLIST_TYPE",
	} {
		t.Setenv(key, "")
	}

	_ = os.Unsetenv("XGG_PROVIDER")
	cmd := newTestRootCmd()
	cfg := resolveRuntimeConfig(cmd)
	def := api.DefaultRuntimeConfig()
	if cfg != def {
		t.Fatalf("expected defaults %+v, got %+v", def, cfg)
	}
}
