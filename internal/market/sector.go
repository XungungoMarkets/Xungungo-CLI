package market

// SectorSummary holds the aggregated % change for a market sector.
type SectorSummary struct {
	Sector    string  `json:"sector"`
	AvgChange float64 `json:"avg_change_pct"`
	Count     int     `json:"count"`
}

// IndustrySummary holds the aggregated % change for a sector/industry pair.
type IndustrySummary struct {
	Sector    string  `json:"sector"`
	Industry  string  `json:"industry"`
	AvgChange float64 `json:"avg_change_pct"`
	Count     int     `json:"count"`
}
