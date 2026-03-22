package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/chart"
	"github.com/XungungoMarkets/xgg/internal/provider"
	"github.com/spf13/cobra"
)

var (
	chartPeriod     string
	chartType       string
	chartOutput     string
	chartWidth      int
	chartHeight     int
	chartTheme      string
	chartInterval   string
	chartIndicators string
)

var chartCmd = &cobra.Command{
	Use:   "chart [symbol]",
	Short: "Generate a price chart as a PNG image",
	Long:  "Fetch historical data and render a line or candlestick chart saved as a PNG file.",
	Example: `  xgg chart AAPL
  xgg chart NVDA --type candlestick --period 3m
  xgg chart TSLA --period 1y --interval week --indicator sma20,sma50
  xgg chart NVDA --type candlestick --period 6m --indicator bb
  xgg chart MSFT --period 5y --interval month --type candlestick --indicator sma20,ema50
  xgg chart AAPL --output /tmp/aapl.png --theme light`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		symbol := strings.ToUpper(strings.TrimSpace(args[0]))

		interval := chart.ParseInterval(chartInterval)

		// Auto-extend period if indicators need more bars than requested.
		period := chartPeriod
		if minP := chart.MinPeriodForIndicators(chartIndicators, interval); minP != "" &&
			chart.PeriodDays(period) < chart.PeriodDays(minP) {
			fmt.Fprintf(cmd.ErrOrStderr(),
				"Note: period extended to %s to compute %s on %s bars\n",
				minP, chartIndicators, chartInterval)
			period = minP
		}

		bars, meta, err := provider.ServiceHandle().GetHistory(context.Background(), symbol, period)
		if err != nil {
			return fmt.Errorf("error fetching history for %s: %w", symbol, err)
		}
		if meta.FallbackUsed {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: falling back to %s for %s history (%v)\n",
				meta.ProviderUsed, symbol, meta.PrimaryErr)
		}
		if len(bars) == 0 {
			return fmt.Errorf("no historical data available for %s", symbol)
		}

		outPath := chartOutput
		if outPath == "" {
			outPath = fmt.Sprintf("%s_chart.png", strings.ToLower(symbol))
		}

		var data []byte
		switch strings.ToLower(chartType) {
		case "candlestick", "candle", "ohlc":
			data, err = chart.RenderCandlestick(symbol, bars, chartWidth, chartHeight, chartTheme, chartIndicators, interval)
		default:
			data, err = chart.RenderLine(symbol, bars, chartWidth, chartHeight, chartTheme, chartIndicators, interval)
		}
		if err != nil {
			return fmt.Errorf("error generating chart: %w", err)
		}

		if err := os.WriteFile(outPath, data, 0644); err != nil {
			return fmt.Errorf("error writing chart to %s: %w", outPath, err)
		}

		fmt.Println(outPath)
		return nil
	},
}

func init() {
	chartCmd.Flags().StringVarP(&chartPeriod, "period", "p", "1m",
		"Time period: 1w, 2w, 5d, 1m, 2m, 3m, 6m, 9m, 1y, 2y, 3y, 5y, 10y, max")
	chartCmd.Flags().StringVarP(&chartType, "type", "t", "line",
		"Chart type: line, candlestick")
	chartCmd.Flags().StringVarP(&chartInterval, "interval", "i", "day",
		"Bar interval: day, week, month")
	chartCmd.Flags().StringVarP(&chartOutput, "output", "o", "",
		"Output file path (default: <symbol>_chart.png)")
	chartCmd.Flags().IntVar(&chartWidth, "width", 900, "Chart width in pixels")
	chartCmd.Flags().IntVar(&chartHeight, "height", 500, "Chart height in pixels")
	chartCmd.Flags().StringVar(&chartTheme, "theme", "dark",
		"Color theme: light, dark, vivid-light, vivid-dark, ant, grafana")
	chartCmd.Flags().StringVar(&chartIndicators, "indicator", "",
		"Overlay indicators (comma-separated): sma20, sma40, sma50, sma100, sma200, ema12, ema26, ema50, ema100, ema200, bb, linear, cubic")
	rootCmd.AddCommand(chartCmd)
}
