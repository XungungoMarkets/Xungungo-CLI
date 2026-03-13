package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

func PrintHistory(symbol string, bars []market.Bar) {
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
