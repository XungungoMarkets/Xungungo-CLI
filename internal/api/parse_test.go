package api

import "testing"

func TestParseFloat(t *testing.T) {
	tests := []struct {
		in   string
		want float64
		ok   bool
	}{
		{in: "$123.45", want: 123.45, ok: true},
		{in: "+1.2%", want: 1.2, ok: true},
		{in: "1,234,567", want: 1234567, ok: true},
		{in: "N/A", ok: false},
		{in: "", ok: false},
	}

	for _, tt := range tests {
		got, err := parseFloat(tt.in)
		if tt.ok && err != nil {
			t.Fatalf("parseFloat(%q) error = %v", tt.in, err)
		}
		if !tt.ok && err == nil {
			t.Fatalf("parseFloat(%q) expected error", tt.in)
		}
		if tt.ok && got != tt.want {
			t.Fatalf("parseFloat(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		in   string
		want int64
	}{
		{in: "1,234", want: 1234},
		{in: "1.5M", want: 1500000},
		{in: "2B", want: 2000000000},
		{in: "3.1T", want: 3100000000000},
		{in: "N/A", want: 0},
		{in: "", want: 0},
	}

	for _, tt := range tests {
		got, err := parseInt64(tt.in)
		if err != nil {
			t.Fatalf("parseInt64(%q) error = %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("parseInt64(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
