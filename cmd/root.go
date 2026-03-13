package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/XungungoMarkets/xgg/internal/api"
	"github.com/spf13/cobra"
)

// Version is set at build time using -ldflags
var Version = "dev"

// JSONOutput is a global flag to control JSON output format
var JSONOutput bool
var providerMode string
var rateLimit int
var maxRetries int
var retryDelay time.Duration
var requestTimeout time.Duration
var watchlistType string

var rootCmd = &cobra.Command{
	Use:   "xgg",
	Short: "Xungungo CLI - Financial markets at your fingertips",
	Long:  "Xungungo CLI provides real-time stock quotes, historical data, and portfolio tracking from your terminal.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set JSON output flag based on parent command
		if cmd.Parent() != nil {
			if JSONOutput, _ = cmd.Parent().Flags().GetBool("json"); !JSONOutput {
				JSONOutput, _ = cmd.Flags().GetBool("json")
			}
		}

		api.ConfigureRuntime(resolveRuntimeConfig(cmd))
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&JSONOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&providerMode, "provider", "", "Data provider: auto, nasdaq, legacy")
	rootCmd.PersistentFlags().IntVar(&rateLimit, "rate-limit", 0, "Max requests per second for Nasdaq provider")
	rootCmd.PersistentFlags().IntVar(&maxRetries, "max-retries", -1, "Max retries for Nasdaq requests")
	rootCmd.PersistentFlags().DurationVar(&retryDelay, "retry-delay", 0, "Retry delay for Nasdaq requests (e.g. 2s)")
	rootCmd.PersistentFlags().DurationVar(&requestTimeout, "timeout", 0, "HTTP timeout for Nasdaq requests (e.g. 30s)")
	rootCmd.PersistentFlags().StringVar(&watchlistType, "watchlist-type", "", "NASDAQ watchlist type parameter (e.g. Rv)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func resolveRuntimeConfig(cmd *cobra.Command) api.RuntimeConfig {
	cfg := api.DefaultRuntimeConfig()
	flags := cmd.Root().PersistentFlags()

	if v := strings.TrimSpace(os.Getenv("XGG_PROVIDER")); v != "" {
		cfg.Provider = api.ProviderMode(strings.ToLower(v))
	}
	if v := strings.TrimSpace(os.Getenv("XGG_RATE_LIMIT")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			cfg.RateLimit = parsed
		}
	}
	if v := strings.TrimSpace(os.Getenv("XGG_MAX_RETRIES")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			cfg.MaxRetries = parsed
		}
	}
	if v := strings.TrimSpace(os.Getenv("XGG_RETRY_DELAY")); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			cfg.RetryDelaySec = int(parsed / time.Second)
		}
	}
	if v := strings.TrimSpace(os.Getenv("XGG_TIMEOUT")); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			cfg.TimeoutSec = int(parsed / time.Second)
		}
	}
	if v := strings.TrimSpace(os.Getenv("XGG_WATCHLIST_TYPE")); v != "" {
		cfg.WatchlistType = v
	}

	if flags.Changed("provider") {
		cfg.Provider = api.ProviderMode(strings.ToLower(strings.TrimSpace(providerMode)))
	}
	if flags.Changed("rate-limit") {
		cfg.RateLimit = rateLimit
	}
	if flags.Changed("max-retries") {
		cfg.MaxRetries = maxRetries
	}
	if flags.Changed("retry-delay") {
		cfg.RetryDelaySec = int(retryDelay / time.Second)
	}
	if flags.Changed("timeout") {
		cfg.TimeoutSec = int(requestTimeout / time.Second)
	}
	if flags.Changed("watchlist-type") {
		cfg.WatchlistType = strings.TrimSpace(watchlistType)
	}

	return cfg
}
