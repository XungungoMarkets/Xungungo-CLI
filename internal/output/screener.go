package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

// Screener column widths (visible characters, excluding ANSI codes).
const (
	scrSymCol  = 6
	scrNameCol = 22
	scrPriceCol = 9
	scrChgCol  = 9
	scrPctCol  = 9
	scrVolCol  = 12
	scrCapCol  = 12
	scrCntryCol = 8
	scrIPOCol  = 4
	scrSecCol  = 18
	scrIndCol  = 20

	// lineWidth = 2 + scrSymCol + 2 + scrNameCol + 2 + scrPriceCol + 2 + scrChgCol + 2 +
	//             scrPctCol + 2 + scrVolCol + 2 + scrCapCol + 2 + scrCntryCol + 2 +
	//             scrIPOCol + 2 + scrSecCol + 2 + scrIndCol + 2
	// = 2 + (6+22+9+9+9+12+12+8+4+18+20) + 11*2+2 = 2 + 129 + 24 = 155
	scrLineWidth = 155
)

func PrintScreenerRows(rows []market.ScreenerRow) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	border := strings.Repeat("─", scrLineWidth)
	printRow := func(line string) {
		fmt.Printf("%s%s%s\n", dim("│"), line, dim("│"))
	}
	printBorder := func(l, r string) {
		fmt.Printf("%s%s%s\n", dim(l), dim(border), dim(r))
	}

	// Header: pad plain text first, then bold.
	hdr := fmt.Sprintf("  %s  %s  %s  %s  %s  %s  %s  %s  %s  %s  %s  ",
		bold(fmt.Sprintf("%-*s", scrSymCol, "SYM")),
		bold(fmt.Sprintf("%-*s", scrNameCol, "Name")),
		bold(fmt.Sprintf("%*s", scrPriceCol, "Price")),
		bold(fmt.Sprintf("%*s", scrChgCol, "Change")),
		bold(fmt.Sprintf("%*s", scrPctCol, "Chg%")),
		bold(fmt.Sprintf("%*s", scrVolCol, "Volume")),
		bold(fmt.Sprintf("%*s", scrCapCol, "MktCap")),
		bold(fmt.Sprintf("%-*s", scrCntryCol, "Country")),
		bold(fmt.Sprintf("%-*s", scrIPOCol, "IPO")),
		bold(fmt.Sprintf("%-*s", scrSecCol, "Sector")),
		bold(fmt.Sprintf("%-*s", scrIndCol, "Industry")),
	)

	printBorder("┌", "┐")
	printRow(hdr)
	printBorder("├", "┤")

	for _, r := range rows {
		pctPadded := fmt.Sprintf("%*s", scrPctCol, truncate(r.PercentageChange, scrPctCol))
		var pctColored string
		switch {
		case strings.HasPrefix(r.PercentageChange, "-"):
			pctColored = red(pctPadded)
		case r.PercentageChange != "" && r.PercentageChange != "N/A":
			pctColored = green(pctPadded)
		default:
			pctColored = pctPadded
		}

		row := fmt.Sprintf("  %-*s  %-*s  %*s  %*s  %s  %*s  %*s  %-*s  %-*s  %-*s  %-*s  ",
			scrSymCol, truncate(r.Symbol, scrSymCol),
			scrNameCol, truncate(r.Name, scrNameCol),
			scrPriceCol, truncate(r.LastSalePrice, scrPriceCol),
			scrChgCol, truncate(r.NetChange, scrChgCol),
			pctColored,
			scrVolCol, truncate(r.Volume, scrVolCol),
			scrCapCol, truncate(r.MarketCap, scrCapCol),
			scrCntryCol, truncate(r.Country, scrCntryCol),
			scrIPOCol, truncate(r.IPOYear, scrIPOCol),
			scrSecCol, truncate(r.Sector, scrSecCol),
			scrIndCol, truncate(r.Industry, scrIndCol),
		)
		printRow(row)
	}

	printBorder("└", "┘")
	fmt.Printf("  %d stocks\n", len(rows))
}
