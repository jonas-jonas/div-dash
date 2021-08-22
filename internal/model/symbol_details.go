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

type SymbolIndicator struct {
	Label  string  `json:"label"`
	Format string  `json:"format"`
	Value  float64 `json:"value"`
}

type SymbolDetails struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Tags        []SymbolTag       `json:"tags"`
	Indicators  []SymbolIndicator `json:"indicators"`
	Description string            `json:"description"`
	Dates       []SymbolDate      `json:"dates"`
}
