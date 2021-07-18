package model

type Chart = []ChartEntry

type ChartEntry struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}
