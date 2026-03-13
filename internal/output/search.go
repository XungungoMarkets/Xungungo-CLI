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
