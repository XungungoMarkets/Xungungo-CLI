package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

func PrintQuote(q *market.StockQuote) {
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
