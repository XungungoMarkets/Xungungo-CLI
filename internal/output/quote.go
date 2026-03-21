package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

func coloredChange(green, red func(...interface{}) string, change, pct float64) (arrow, text string) {
	if change >= 0 {
		return green("▲"), green(fmt.Sprintf("+%.2f (+%.2f%%)", change, pct))
	}
	return red("▼"), red(fmt.Sprintf("%.2f (%.2f%%)", change, pct))
}

func PrintQuote(q *market.StockQuote) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Market state badge
	var stateBadge string
	switch q.MarketState {
	case "REGULAR":
		stateBadge = green("● Open")
	case "PRE", "PREPRE":
		stateBadge = yellow("● Pre-Market")
	case "POST", "POSTPOST":
		stateBadge = yellow("● After Hours")
	case "CLOSED":
		stateBadge = dim("● Closed")
	}

	name := q.Name
	if name == "" {
		name = q.Symbol
	}

	// Build lines (colored)
	var lines []string

	// Line 1: symbol - name  [badge]
	line1 := fmt.Sprintf("%s - %s", bold(q.Symbol), name)
	if stateBadge != "" {
		line1 = fmt.Sprintf("%s - %s   %s", bold(q.Symbol), name, stateBadge)
	}
	lines = append(lines, line1)

	// Line 2: main price — context-aware label during extended sessions
	arrow, changeStr := coloredChange(green, red, q.Change, q.ChangePercent)
	priceFmt := bold(fmt.Sprintf("$%.2f", q.Price))
	switch q.MarketState {
	case "POST", "POSTPOST":
		// Primary data is the regular-session close; show it labeled with day change
		lines = append(lines, fmt.Sprintf("%s  %s  %s %s", dim("Close:"), priceFmt, arrow, changeStr))
	case "PRE", "PREPRE":
		// Primary data is yesterday's close; show price only (change belongs to pre-market line)
		lines = append(lines, fmt.Sprintf("%s  %s", dim("Prev Close:"), priceFmt))
	default:
		lines = append(lines, fmt.Sprintf("%s  %s %s", priceFmt, arrow, changeStr))
	}

	// Previous close (from providers that supply it separately, e.g. Yahoo)
	if q.PreviousClose > 0 {
		lines = append(lines, fmt.Sprintf("%s  $%.2f", dim("Prev Close:"), q.PreviousClose))
	}

	// Pre-market session
	if q.PreMarketPrice > 0 {
		pmArrow, pmChange := coloredChange(green, red, q.PreMarketChange, q.PreMarketChangePercent)
		lines = append(lines, fmt.Sprintf("%s  %s  %s %s", dim("Pre-Market:"), bold(fmt.Sprintf("$%.2f", q.PreMarketPrice)), pmArrow, pmChange))
	}

	// After-hours session
	if q.PostMarketPrice > 0 {
		ahArrow, ahChange := coloredChange(green, red, q.PostMarketChange, q.PostMarketChangePercent)
		lines = append(lines, fmt.Sprintf("%s  %s  %s %s", dim("After Hours:"), bold(fmt.Sprintf("$%.2f", q.PostMarketPrice)), ahArrow, ahChange))
	}

	// Last line: volume + market cap
	lines = append(lines, fmt.Sprintf("Vol: %s  │  Mkt Cap: %s", formatNumber(q.Volume), formatCap(q.MarketCap)))

	// Compute content width from visible lengths
	contentWidth := 0
	for _, l := range lines {
		if vl := visLen(l); vl > contentWidth {
			contentWidth = vl
		}
	}
	boxInner := contentWidth + 4 // 2-space margin on each side

	border := dim(strings.Repeat("─", boxInner))
	printRow := func(line string) {
		fmt.Printf("%s  %s  %s\n", dim("│"), padRight(line, contentWidth), dim("│"))
	}

	fmt.Printf("%s%s%s\n", dim("┌"), border, dim("┐"))
	for i, l := range lines {
		if i == 1 {
			// separator after title
			fmt.Printf("%s%s%s\n", dim("├"), border, dim("┤"))
		}
		printRow(l)
	}
	fmt.Printf("%s%s%s\n", dim("└"), border, dim("┘"))
}
