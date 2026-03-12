package api

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func parseFloat(input string) (float64, error) {
	clean := normalizeNumber(input)
	if clean == "" {
		return 0, fmt.Errorf("empty numeric value")
	}
	v, err := strconv.ParseFloat(clean, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("invalid float %q", input)
	}
	return v, nil
}

func parseInt(input string) (int, error) {
	v, err := parseInt64(input)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func parseInt64(input string) (int64, error) {
	raw := strings.TrimSpace(input)
	if raw == "" || strings.EqualFold(raw, "n/a") || raw == "-" {
		return 0, nil
	}

	multiplier := float64(1)
	last := raw[len(raw)-1]
	switch last {
	case 'K', 'k':
		multiplier = 1_000
		raw = raw[:len(raw)-1]
	case 'M', 'm':
		multiplier = 1_000_000
		raw = raw[:len(raw)-1]
	case 'B', 'b':
		multiplier = 1_000_000_000
		raw = raw[:len(raw)-1]
	case 'T', 't':
		multiplier = 1_000_000_000_000
		raw = raw[:len(raw)-1]
	}

	base, err := parseFloat(raw)
	if err != nil {
		return 0, err
	}
	return int64(base * multiplier), nil
}

func normalizeNumber(input string) string {
	s := strings.TrimSpace(input)
	if s == "" || strings.EqualFold(s, "n/a") || s == "-" {
		return ""
	}

	s = strings.TrimPrefix(s, "$")
	s = strings.TrimPrefix(s, "+")
	s = strings.TrimSuffix(s, "%")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "(", "-")
	s = strings.ReplaceAll(s, ")", "")
	return strings.TrimSpace(s)
}
