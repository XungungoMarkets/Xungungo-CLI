package output

import (
	"fmt"
	"strings"

	"github.com/XungungoMarkets/xgg/internal/market"
	"github.com/fatih/color"
)

// Column widths (visible characters, excluding ANSI codes).
const (
	indNameCol = 52 // name column: sector uses full width, industry gets 2-space indent
	indPctCol  = 8  // pct column, right-aligned
	indCntCol  = 6  // count column, right-aligned
	// Total visible content per row: 2(│pad) + 1(▲) + 1( ) + 52 + 2 + 8 + 2 + 6 + 2(│pad) = 76
	indLineWidth = 2 + 1 + 1 + indNameCol + 2 + indPctCol + 2 + indCntCol + 2
)

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// fmtPctField returns a right-aligned pct string of fixed width (plain text, no ANSI).
func fmtPctField(v float64) string {
	var raw string
	if v >= 0 {
		raw = fmt.Sprintf("+%.2f%%", v)
	} else {
		raw = fmt.Sprintf("%.2f%%", v)
	}
	return fmt.Sprintf("%*s", indPctCol, raw)
}

// sectorGroup groups industries under a sector with a pre-computed weighted avg.
type sectorGroup struct {
	sector    string
	avgChange float64
	count     int
	rows      []market.IndustrySummary
}

func buildSectorGroups(industries []market.IndustrySummary) []sectorGroup {
	var groups []sectorGroup
	idx := map[string]int{}
	for _, ind := range industries {
		i, ok := idx[ind.Sector]
		if !ok {
			i = len(groups)
			idx[ind.Sector] = i
			groups = append(groups, sectorGroup{sector: ind.Sector})
		}
		groups[i].rows = append(groups[i].rows, ind)
		groups[i].count += ind.Count
		groups[i].avgChange += ind.AvgChange * float64(ind.Count)
	}
	for i := range groups {
		if groups[i].count > 0 {
			groups[i].avgChange /= float64(groups[i].count)
		}
	}
	// sort sectors by weighted avg change desc
	for i := 1; i < len(groups); i++ {
		for j := i; j > 0 && groups[j].avgChange > groups[j-1].avgChange; j-- {
			groups[j], groups[j-1] = groups[j-1], groups[j]
		}
	}
	return groups
}

func PrintIndustries(industries []market.IndustrySummary) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	border := strings.Repeat("─", indLineWidth)
	// printRow prints a pre-built visible-width line inside │ borders.
	// `line` must be exactly indLineWidth visible chars.
	printRow := func(line string) {
		fmt.Printf("%s%s%s\n", dim("│"), line, dim("│"))
	}
	printBorder := func(l, m, r string) {
		fmt.Printf("%s%s%s\n", dim(l), dim(border), dim(r))
	}

	// Header
	headerName := fmt.Sprintf("%-*s", indNameCol, "Industry")
	header := fmt.Sprintf("  %s %s  %*s  %*s  ",
		" ", bold(headerName), indPctCol, bold("Avg Chg%"), indCntCol, bold("Stocks"))
	printBorder("┌", "─", "┐")
	printRow(header)

	groups := buildSectorGroups(industries)
	for _, g := range groups {
		printBorder("├", "─", "┤")

		// Sector header: bold cyan name, weighted pct
		sectorArrow := green("▲")
		pctField := fmtPctField(g.avgChange)
		coloredPct := green(pctField)
		if g.avgChange < 0 {
			sectorArrow = red("▼")
			coloredPct = red(pctField)
		}
		// Pad plain name first, then colorize the padded field
		sectorNameField := fmt.Sprintf("%-*s", indNameCol, truncate(g.sector, indNameCol))
		printRow(fmt.Sprintf("  %s %s  %s  %*d  ",
			sectorArrow, cyan(bold(sectorNameField)), coloredPct, indCntCol, g.count))

		// Industry rows: 2-space indent, name shrunk by 2
		const indIndent = 2
		const indIndNameCol = indNameCol - indIndent
		for _, ind := range g.rows {
			arrow := green("▲")
			pf := fmtPctField(ind.AvgChange)
			cp := green(pf)
			if ind.AvgChange < 0 {
				arrow = red("▼")
				cp = red(pf)
			}
			// Indent + padded plain name + color
			nameField := fmt.Sprintf("%*s%-*s",
				indIndent, "", indIndNameCol, truncate(ind.Industry, indIndNameCol))
			printRow(fmt.Sprintf("  %s %s  %s  %*d  ",
				arrow, nameField, cp, indCntCol, ind.Count))
		}
	}

	printBorder("└", "─", "┘")
}

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
