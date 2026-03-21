package output

import (
	"fmt"
	"strings"
)

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

// visLen returns the visible (terminal) length of a string, excluding ANSI escape codes.
func visLen(s string) int {
	n := 0
	inEscape := false
	for _, r := range s {
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if r == '\x1b' {
			inEscape = true
			continue
		}
		n++
	}
	return n
}

// padRight pads s with spaces on the right to reach the given visible width.
func padRight(s string, width int) string {
	pad := width - visLen(s)
	if pad <= 0 {
		return s
	}
	return s + strings.Repeat(" ", pad)
}
