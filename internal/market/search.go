package market

// SearchResult represents a symbol search suggestion.
type SearchResult struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
