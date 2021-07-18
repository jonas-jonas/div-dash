package model

type SymbolTag struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	Link  string `json:"link"`
}

type SymbolDate struct {
	Label string `json:"label"`
	Date  string `json:"date"`
}

type SymbolDetails struct {
	Type          string       `json:"type"`
	Name          string       `json:"name"`
	Tags          []SymbolTag  `json:"tags"`
	MarketCap     int64        `json:"marketCap"`
	PERatio       float64      `json:"peRatio"`
	DividendYield float64      `json:"dividendYield"`
	EPS           float64      `json:"eps"`
	Description   string       `json:"description"`
	Dates         []SymbolDate `json:"dates"`
}
