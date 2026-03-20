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

// CountrySummary holds the aggregated % change for a country.
type CountrySummary struct {
	Country   string  `json:"country"`
	AvgChange float64 `json:"avg_change_pct"`
	Count     int     `json:"count"`
}

// StockDetail holds per-stock data for display in grouped views.
type StockDetail struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	ChangePct float64 `json:"change_pct"`
}

// CountryWithStocks holds a country summary plus its individual stocks.
type CountryWithStocks struct {
	Country   string        `json:"country"`
	AvgChange float64       `json:"avg_change_pct"`
	Count     int           `json:"count"`
	Stocks    []StockDetail `json:"stocks"`
}
