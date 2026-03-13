package chart

import (
	"fmt"
	"strings"

	"github.com/go-analyze/charts"
)

// IndicatorDef describes a named chart indicator as one or more trend lines.
// Adding a new indicator is as simple as adding a new entry to the Registry.
type IndicatorDef struct {
	// Label is shown in legend / error messages.
	Label string
	// TrendLines are the underlying chart trend lines to render.
	TrendLines []charts.SeriesTrendLine
}

// Registry maps indicator names (lower-case) to their definitions.
// Extend this map to add new indicators — no other code needs to change.
var Registry = map[string]IndicatorDef{
	// Simple Moving Averages
	"sma20": {
		Label: "SMA(20)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeSMA, Period: 20},
		},
	},
	"sma50": {
		Label: "SMA(50)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeSMA, Period: 50},
		},
	},
	"sma200": {
		Label: "SMA(200)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeSMA, Period: 200},
		},
	},

	// Exponential Moving Averages
	"ema12": {
		Label: "EMA(12)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeEMA, Period: 12},
		},
	},
	"ema26": {
		Label: "EMA(26)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeEMA, Period: 26},
		},
	},
	"ema50": {
		Label: "EMA(50)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeEMA, Period: 50},
		},
	},

	// Bollinger Bands (upper + middle SMA + lower)
	"bb": {
		Label: "BB(20)",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeBollingerUpper, Period: 20},
			{Type: charts.SeriesTrendTypeSMA, Period: 20},
			{Type: charts.SeriesTrendTypeBollingerLower, Period: 20},
		},
	},

	// Trend lines
	"linear": {
		Label: "Linear",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeLinear},
		},
	},
	"cubic": {
		Label: "Cubic",
		TrendLines: []charts.SeriesTrendLine{
			{Type: charts.SeriesTrendTypeCubic},
		},
	},
}

// ResolveIndicators parses a comma-separated indicator string and returns
// the combined list of trend lines plus a legend label list.
// Returns an error if any unknown indicator name is used.
func ResolveIndicators(raw string) ([]charts.SeriesTrendLine, []string, error) {
	if raw == "" {
		return nil, nil, nil
	}

	var trendLines []charts.SeriesTrendLine
	var labels []string

	for _, name := range strings.Split(raw, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		def, ok := Registry[name]
		if !ok {
			return nil, nil, fmt.Errorf("unknown indicator %q — available: %s", name, AvailableIndicators())
		}
		trendLines = append(trendLines, def.TrendLines...)
		labels = append(labels, def.Label)
	}

	return trendLines, labels, nil
}

// AvailableIndicators returns a sorted, comma-separated list of registered indicator names.
func AvailableIndicators() string {
	names := make([]string, 0, len(Registry))
	for k := range Registry {
		names = append(names, k)
	}
	// deterministic order
	sorted := sortedKeys(names)
	return strings.Join(sorted, ", ")
}

func sortedKeys(keys []string) []string {
	// insertion sort — small slice, no need for sort package import
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
