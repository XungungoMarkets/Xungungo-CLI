package chart

import (
	"fmt"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/go-analyze/charts"
)

const (
	defaultLabelCount = 8
)

// RenderLine generates a PNG line chart of closing prices for the given bars.
func RenderLine(symbol string, bars []market.Bar, width, height int, theme string) ([]byte, error) {
	if len(bars) == 0 {
		return nil, fmt.Errorf("no data to render")
	}

	closes := make([]float64, len(bars))
	labels := make([]string, len(bars))
	for i, b := range bars {
		closes[i] = b.Close
		labels[i] = b.Date.Format("01/02")
	}

	opt := charts.NewLineChartOptionWithData([][]float64{closes})
	opt.Title.Text = fmt.Sprintf("%s — Close Price", symbol)
	opt.Title.FontStyle.FontSize = 14
	opt.XAxis.Labels = labels
	opt.XAxis.LabelCount = labelCount(len(bars))
	opt.Legend.SeriesNames = []string{"Close"}
	opt.Theme = charts.GetTheme(theme)
	opt.LineStrokeWidth = 1.5
	opt.Symbol = charts.SymbolDot
	opt.FillArea = charts.Ptr(true)
	opt.FillOpacity = 40
	opt.Padding = charts.NewBoxEqual(20)

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

// RenderCandlestick generates a PNG candlestick (OHLC) chart for the given bars.
func RenderCandlestick(symbol string, bars []market.Bar, width, height int, theme string) ([]byte, error) {
	if len(bars) == 0 {
		return nil, fmt.Errorf("no data to render")
	}

	ohlc := make([]charts.OHLCData, len(bars))
	labels := make([]string, len(bars))
	for i, b := range bars {
		ohlc[i] = charts.OHLCData{Open: b.Open, High: b.High, Low: b.Low, Close: b.Close}
		labels[i] = b.Date.Format("01/02")
	}

	opt := charts.CandlestickChartOption{
		Theme: charts.GetTheme(theme),
		SeriesList: charts.CandlestickSeriesList{
			{Data: ohlc, Name: symbol},
		},
		XAxis: charts.XAxisOption{
			Labels:     labels,
			LabelCount: labelCount(len(bars)),
		},
		Legend: charts.LegendOption{
			SeriesNames: []string{symbol},
		},
		Padding: charts.NewBoxEqual(20),
	}
	opt.Title.Text = fmt.Sprintf("%s — Candlestick", symbol)
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

// labelCount returns a reasonable number of X-axis labels for n data points.
func labelCount(n int) int {
	if n <= defaultLabelCount {
		return n
	}
	return defaultLabelCount
}
