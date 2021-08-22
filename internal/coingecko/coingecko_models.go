package coingecko

type CoingeckoCoin struct {
	CoingeckoID string `json:"id"`
	SymbolID    string `json:"symbol"`
	Name        string `json:"name"`
}

// Example:
// {
//     "polkadot": {
//         "eur": 23.7
//     }
// }
type CoingeckoPriceResponse map[string]map[string]float64
