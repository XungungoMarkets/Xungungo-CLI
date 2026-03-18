package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

func PrintSectors(sectors []market.SectorSummary) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	header := fmt.Sprintf("  %-35s  %8s  %6s", "Sector", "Avg Chg%", "Stocks")
	width := len(header) + 2
	border := dim(strings.Repeat("─", width))

	fmt.Printf("%s%s%s\n", dim("┌"), border, dim("┐"))
	fmt.Printf("%s  %-*s  %s\n", dim("│"), width-2, bold(header), dim("│"))
	fmt.Printf("%s%s%s\n", dim("├"), border, dim("┤"))

	for _, s := range sectors {
		arrow := green("▲")
		pctStr := green(fmt.Sprintf("+%.2f%%", s.AvgChange))
		if s.AvgChange < 0 {
			arrow = red("▼")
			pctStr = red(fmt.Sprintf("%.2f%%", s.AvgChange))
		}
		row := fmt.Sprintf("  %s %-33s  %8s  %6d", arrow, s.Sector, pctStr, s.Count)
		fmt.Printf("%s  %-*s  %s\n", dim("│"), width-2, row, dim("│"))
	}

	fmt.Printf("%s%s%s\n", dim("└"), border, dim("┘"))
}
