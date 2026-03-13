package market

import "fmt"

// ApplyInterval aggregates bars by the given interval string ("week", "month").
// Any other value (including "day") returns bars unchanged.
func ApplyInterval(bars []Bar, interval string) []Bar {
	switch interval {
	case "week", "w", "weekly":
		return AggregateWeekly(bars)
	case "month", "m", "monthly":
		return AggregateMonthly(bars)
	default:
		return bars
	}
}

// AggregateWeekly groups daily bars into weekly bars.
// Each week runs Mon→Fri: Open from first day, Close from last day, High/Low extremes, Volume summed.
func AggregateWeekly(bars []Bar) []Bar {
	return aggregate(bars, func(b Bar) string {
		y, w := b.Date.ISOWeek()
		return formatKey(y, w)
	})
}

// AggregateMonthly groups daily bars into monthly bars.
func AggregateMonthly(bars []Bar) []Bar {
	return aggregate(bars, func(b Bar) string {
		return b.Date.Format("2006-01")
	})
}

func aggregate(bars []Bar, key func(Bar) string) []Bar {
	if len(bars) == 0 {
		return nil
	}

	type bucket struct {
		bar Bar
	}

	seen := make(map[string]*bucket)
	order := make([]string, 0)

	for _, b := range bars {
		k := key(b)
		if existing, ok := seen[k]; !ok {
			seen[k] = &bucket{bar: b}
			order = append(order, k)
		} else {
			if b.High > existing.bar.High {
				existing.bar.High = b.High
			}
			if b.Low < existing.bar.Low {
				existing.bar.Low = b.Low
			}
			existing.bar.Close = b.Close
			existing.bar.Volume += b.Volume
		}
	}

	result := make([]Bar, 0, len(order))
	for _, k := range order {
		result = append(result, seen[k].bar)
	}
	return result
}

func formatKey(year, week int) string {
	// zero-pad week for correct lexicographic ordering
	if week < 10 {
		return fmt.Sprintf("%d-W0%d", year, week)
	}
	return fmt.Sprintf("%d-W%d", year, week)
}
