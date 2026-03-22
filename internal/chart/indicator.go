package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-analyze/charts"
	"github.com/markcheno/go-talib"
)

// IndicatorDef describes a named chart indicator as one or more trend lines.
// Adding a new indicator is as simple as adding a new entry to the Registry.
type IndicatorDef struct {
	// Label is shown in legend / error messages.
	Label string
	// TrendLines are the underlying chart trend lines to render.
	// Not used when MAType is set.
	TrendLines []charts.SeriesTrendLine
	// MAType is "sma" or "ema" for indicators that must be pre-computed with
	// talib (correct trailing window). When set, TrendLines is ignored and
	// the result is injected as a separate data series.
	MAType   string
	MAPeriod int
}

// Registry maps indicator names (lower-case) to their definitions.
// Extend this map to add new indicators — no other code needs to change.
var Registry = map[string]IndicatorDef{
	// Simple Moving Averages — pre-computed with talib (correct trailing window)
	"sma20":  {Label: "SMA(20)", MAType: "sma", MAPeriod: 20},
	"sma40":  {Label: "SMA(40)", MAType: "sma", MAPeriod: 40},
	"sma50":  {Label: "SMA(50)", MAType: "sma", MAPeriod: 50},
	"sma100": {Label: "SMA(100)", MAType: "sma", MAPeriod: 100},
	"sma200": {Label: "SMA(200)", MAType: "sma", MAPeriod: 200},

	// Exponential Moving Averages — pre-computed with talib
	"ema12":  {Label: "EMA(12)", MAType: "ema", MAPeriod: 12},
	"ema26":  {Label: "EMA(26)", MAType: "ema", MAPeriod: 26},
	"ema50":  {Label: "EMA(50)", MAType: "ema", MAPeriod: 50},
	"ema100": {Label: "EMA(100)", MAType: "ema", MAPeriod: 100},
	"ema200": {Label: "EMA(200)", MAType: "ema", MAPeriod: 200},

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

// ResolvedIndicators holds the results of parsing an indicator string.
type ResolvedIndicators struct {
	// TrendLines are passed directly to the chart library (bb, linear, cubic…).
	TrendLines []charts.SeriesTrendLine
	// PrecomputedSeries are trailing MA series computed with talib.
	// Each entry is a []float64 slice aligned to the closes array.
	PrecomputedSeries [][]float64
	// Labels contains one label per indicator (all types), in order.
	Labels []string
}

// ResolveIndicators parses a comma-separated indicator string.
// MA indicators (sma*, ema*) are marked for talib pre-computation;
// others (bb, linear, cubic) use the chart library's built-in trend lines.
// closes must be provided for pre-computation; pass nil to skip it (e.g. for validation).
func ResolveIndicators(raw string, closes []float64) (*ResolvedIndicators, error) {
	res := &ResolvedIndicators{}
	if raw == "" {
		return res, nil
	}

	for _, name := range strings.Split(raw, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		def, ok := Registry[name]
		if !ok {
			return nil, fmt.Errorf("unknown indicator %q — available: %s", name, AvailableIndicators())
		}
		res.Labels = append(res.Labels, def.Label)

		if def.MAType != "" {
			if closes != nil {
				res.PrecomputedSeries = append(res.PrecomputedSeries, ComputeMA(closes, def.MAType, def.MAPeriod))
			}
		} else {
			res.TrendLines = append(res.TrendLines, def.TrendLines...)
		}
	}
	return res, nil
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

// ComputeMA computes a trailing SMA or EMA using talib.
// The first (period-1) values are set to math.MaxFloat64 (library null) so
// they are skipped during rendering.
func ComputeMA(closes []float64, maType string, period int) []float64 {
	var raw []float64
	switch maType {
	case "ema":
		raw = talib.Ema(closes, period)
	default:
		raw = talib.Sma(closes, period)
	}
	result := make([]float64, len(raw))
	for i, v := range raw {
		if v == 0 && i < period-1 {
			result[i] = math.MaxFloat64 // library null value — skipped in rendering
		} else {
			result[i] = v
		}
	}
	return result
}

// CandlestickIndicators returns all indicators as CloseTrendLine entries for candlestick charts.
// MA indicators use the chart library's built-in SMA/EMA (centered moving average).
func CandlestickIndicators(raw string) ([]charts.SeriesTrendLine, []string, error) {
	var trendLines []charts.SeriesTrendLine
	var labels []string
	if raw == "" {
		return trendLines, labels, nil
	}
	for _, name := range strings.Split(raw, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		def, ok := Registry[name]
		if !ok {
			return nil, nil, fmt.Errorf("unknown indicator %q — available: %s", name, AvailableIndicators())
		}
		labels = append(labels, def.Label)
		if def.MAType != "" {
			t := charts.SeriesTrendTypeSMA
			if def.MAType == "ema" {
				t = charts.SeriesTrendTypeEMA
			}
			trendLines = append(trendLines, charts.SeriesTrendLine{Type: t, Period: def.MAPeriod})
		} else {
			trendLines = append(trendLines, def.TrendLines...)
		}
	}
	return trendLines, labels, nil
}

// maxIndicatorPeriod returns the largest period across all requested indicators.
func maxIndicatorPeriod(raw string) int {
	max := 0
	for _, name := range strings.Split(raw, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if def, ok := Registry[name]; ok {
			for _, tl := range def.TrendLines {
				if tl.Period > max {
					max = tl.Period
				}
			}
		}
	}
	return max
}

// periodDays maps a period string to approximate calendar days.
func PeriodDays(p string) int {
	switch p {
	case "5d":
		return 5
	case "1w":
		return 7
	case "2w":
		return 14
	case "1m":
		return 30
	case "2m":
		return 60
	case "3m":
		return 90
	case "6m":
		return 180
	case "9m":
		return 270
	case "1y":
		return 365
	case "2y":
		return 730
	case "3y":
		return 1095
	case "5y":
		return 1825
	case "10y":
		return 3650
	case "max":
		return 7300
	default:
		return 30
	}
}

// MinPeriodForIndicators returns the minimum fetch period needed to have enough
// bars to compute the requested indicators for the given interval.
// Returns "" when no adjustment is needed.
func MinPeriodForIndicators(raw string, interval Interval) string {
	maxP := maxIndicatorPeriod(raw)
	if maxP == 0 {
		return ""
	}

	// Calendar days needed (with a 20% buffer so the visible range still shows
	// meaningful data after the warm-up bars are consumed).
	var calDays int
	switch interval {
	case IntervalWeek:
		calDays = int(float64(maxP*7) * 1.2)
	case IntervalMonth:
		calDays = int(float64(maxP*31) * 1.2)
	default: // day — account for weekends/holidays
		calDays = int(float64(maxP) * 1.5)
	}

	for _, p := range []struct {
		s string
		d int
	}{
		{"1m", 30}, {"2m", 60}, {"3m", 90}, {"6m", 180}, {"9m", 270},
		{"1y", 365}, {"2y", 730}, {"3y", 1095}, {"5y", 1825},
		{"10y", 3650}, {"max", 7300},
	} {
		if p.d >= calDays {
			return p.s
		}
	}
	return "max"
}
