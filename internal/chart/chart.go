package chart

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/go-analyze/charts"
)

// Interval controls how daily bars are aggregated before rendering.
type Interval string

const (
	IntervalDay   Interval = "day"
	IntervalWeek  Interval = "week"
	IntervalMonth Interval = "month"
)

// ParseInterval normalises user input to an Interval constant.
func ParseInterval(s string) Interval {
	switch strings.ToLower(s) {
	case "week", "w", "weekly":
		return IntervalWeek
	case "month", "m", "monthly":
		return IntervalMonth
	default:
		return IntervalDay
	}
}

// RenderLine generates a PNG line chart of closing prices.
// indicators is a comma-separated list of overlay indicator names (e.g. "sma20,bb").
func RenderLine(symbol string, bars []market.Bar, width, height int, theme, indicators string, interval Interval) ([]byte, error) {
	bars = applyInterval(bars, interval)
	if len(bars) == 0 {
		return nil, fmt.Errorf("no data to render")
	}

	trendLines, indLabels, err := ResolveIndicators(indicators)
	if err != nil {
		return nil, err
	}

	closes := make([]float64, len(bars))
	labels := make([]string, len(bars))
	dateFmt := dateLabelFormat(interval)
	for i, b := range bars {
		closes[i] = b.Close
		labels[i] = b.Date.Format(dateFmt)
	}

	seriesName := "Close"
	legendNames := append([]string{seriesName}, indLabels...)

	opt := charts.NewLineChartOptionWithData([][]float64{closes})
	opt.Title.Text = chartTitle(symbol, "Close Price", interval, indLabels)
	opt.Title.FontStyle.FontSize = 14
	opt.XAxis.Labels = labels
	opt.XAxis.LabelCount = labelCount(len(bars))
	opt.Legend.SeriesNames = legendNames
	opt.Theme = charts.GetTheme(theme)
	opt.LineStrokeWidth = 1.5
	opt.Symbol = charts.SymbolDot
	opt.FillArea = charts.Ptr(true)
	opt.FillOpacity = 40
	opt.Padding = charts.NewBoxEqual(20)
	// Attach trend lines to the first (only) series
	opt.SeriesList[0].TrendLine = trendLines

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        width,
		Height:       height,
	}, charts.PainterThemeOption(charts.GetTheme(theme)))

	if err := p.LineChart(opt); err != nil {
		return nil, err
	}
	return p.Bytes()
}

// RenderCandlestick generates a PNG candlestick (OHLC) chart.
// indicators is a comma-separated list of overlay indicator names (e.g. "sma20,bb").
func RenderCandlestick(symbol string, bars []market.Bar, width, height int, theme, indicators string, interval Interval) ([]byte, error) {
	bars = applyInterval(bars, interval)
	if len(bars) == 0 {
		return nil, fmt.Errorf("no data to render")
	}

	trendLines, indLabels, err := ResolveIndicators(indicators)
	if err != nil {
		return nil, err
	}

	ohlc := make([]charts.OHLCData, len(bars))
	labels := make([]string, len(bars))
	dateFmt := dateLabelFormat(interval)
	for i, b := range bars {
		ohlc[i] = charts.OHLCData{Open: b.Open, High: b.High, Low: b.Low, Close: b.Close}
		labels[i] = b.Date.Format(dateFmt)
	}

	legendNames := append([]string{symbol}, indLabels...)

	opt := charts.CandlestickChartOption{
		Theme: charts.GetTheme(theme),
		SeriesList: charts.CandlestickSeriesList{
			{Data: ohlc, Name: symbol, CloseTrendLine: trendLines},
		},
		XAxis: charts.XAxisOption{
			Labels:     labels,
			LabelCount: labelCount(len(bars)),
		},
		Legend: charts.LegendOption{
			SeriesNames: legendNames,
		},
		Padding: charts.NewBoxEqual(20),
	}
	opt.Title.Text = chartTitle(symbol, "Candlestick", interval, indLabels)
	opt.Title.FontStyle.FontSize = 14

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        width,
		Height:       height,
	}, charts.PainterThemeOption(charts.GetTheme(theme)))

	if err := p.CandlestickChart(opt); err != nil {
		return nil, err
	}
	return p.Bytes()
}

func applyInterval(bars []market.Bar, interval Interval) []market.Bar {
	return market.ApplyInterval(bars, string(interval))
}

func dateLabelFormat(interval Interval) string {
	switch interval {
	case IntervalMonth:
		return "Jan 06"
	default:
		return "01/02"
	}
}

func chartTitle(symbol, chartType string, interval Interval, indLabels []string) string {
	title := fmt.Sprintf("%s — %s (%s)", symbol, chartType, string(interval))
	if len(indLabels) > 0 {
		title += " · " + strings.Join(indLabels, " · ")
	}
	return title
}

// labelCount returns a reasonable number of X-axis labels for n data points.
func labelCount(n int) int {
	const max = 10
	if n <= max {
		return n
	}
	return max
}
