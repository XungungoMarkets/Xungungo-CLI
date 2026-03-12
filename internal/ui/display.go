package ui

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/api"
	"github.com/fatih/color"
)

func PrintQuote(q *api.StockQuote) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	arrow := green("▲")
	changeStr := green(fmt.Sprintf("+%.2f (+%.2f%%)", q.Change, q.ChangePercent))
	if q.Change < 0 {
		arrow = red("▼")
		changeStr = red(fmt.Sprintf("%.2f (%.2f%%)", q.Change, q.ChangePercent))
	}

	name := q.Name
	if name == "" {
		name = q.Symbol
	}

	line1 := fmt.Sprintf("  %s - %s", bold(q.Symbol), name)
	line2 := fmt.Sprintf("  $%.2f  %s %s", q.Price, arrow, changeStr)
	line3 := fmt.Sprintf("  Vol: %s  │  Mkt Cap: %s", formatNumber(q.Volume), formatCap(q.MarketCap))

	width := maxLen(line1, line2, line3) + 4
	border := dim(strings.Repeat("─", width))

	fmt.Printf("%s%s%s\n", dim("┌"), border, dim("┐"))
	fmt.Printf("%s  %-*s  %s\n", dim("│"), width-2, line1, dim("│"))
	fmt.Printf("%s  %-*s  %s\n", dim("│"), width-2, line2, dim("│"))
	fmt.Printf("%s  %-*s  %s\n", dim("│"), width-2, line3, dim("│"))
	fmt.Printf("%s%s%s\n", dim("└"), border, dim("┘"))
}

func PrintHistory(symbol string, bars []api.Bar) {
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	fmt.Printf("\n  %s Historical Data\n", bold(symbol))
	fmt.Printf("  %s\n", dim("Date          Open      High       Low     Close       Volume"))
	fmt.Printf("  %s\n", dim(strings.Repeat("─", 66)))

	for i, b := range bars {
		date := b.Date.Format("2006-01-02")

		changeMarker := ""
		if i > 0 {
			if b.Close > bars[i-1].Close {
				changeMarker = green(" ▲")
			} else if b.Close < bars[i-1].Close {
				changeMarker = red(" ▼")
			}
		}

		fmt.Printf("  %s  %8.2f  %8.2f  %8.2f  %8.2f  %10s%s\n",
			date, b.Open, b.High, b.Low, b.Close, formatNumber(b.Volume), changeMarker)
	}
	fmt.Println()
}

func PrintSearchResults(query string, results []api.SearchResult) {
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	fmt.Printf("\n  %s Results for %q\n", bold("Search"), query)
	if len(results) == 0 {
		fmt.Printf("  %s\n\n", dim("No results found."))
		return
	}

	fmt.Printf("  %s\n", dim("Symbol    Type      Name"))
	fmt.Printf("  %s\n", dim(strings.Repeat("-", 72)))
	for _, r := range results {
		fmt.Printf("  %-8s  %-8s  %s\n", r.Symbol, r.Type, r.Name)
		if r.Description != "" {
			fmt.Printf("                    %s\n", dim(r.Description))
		}
	}
	fmt.Println()
}

func formatNumber(n int) string {
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

func formatCap(n int64) string {
	if n >= 1_000_000_000_000 {
		return fmt.Sprintf("%.1fT", float64(n)/1_000_000_000_000)
	}
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	return fmt.Sprintf("%d", n)
}

func maxLen(strs ...string) int {
	max := 0
	for _, s := range strs {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}
