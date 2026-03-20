package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

func PrintSearchResults(query string, results []market.SearchResult) {
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("\n  %s Results for %q\n", bold("Search"), query)
	if len(results) == 0 {
		fmt.Printf("  %s\n\n", dim("No results found."))
		return
	}

	fmt.Printf("  %s\n", dim("Symbol    Source    Type      Name"))
	fmt.Printf("  %s\n", dim(strings.Repeat("-", 80)))
	for _, r := range results {
		src := r.Source
		var coloredSrc string
		switch src {
		case "NASDAQ":
			coloredSrc = cyan(fmt.Sprintf("%-8s", src))
		case "Yahoo":
			coloredSrc = yellow(fmt.Sprintf("%-8s", src))
		default:
			coloredSrc = fmt.Sprintf("%-8s", src)
		}
		fmt.Printf("  %-8s  %s  %-8s  %s\n", r.Symbol, coloredSrc, r.Type, r.Name)
		if r.Description != "" {
			fmt.Printf("                              %s\n", dim(r.Description))
		}
	}
	fmt.Println()
}
