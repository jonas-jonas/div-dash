package iex

type CompanyDetails struct {
	Symbol         string   `json:"symbol"`
	CompanyName    string   `json:"companyName"`
	Exchange       string   `json:"exchange"`
	Industry       string   `json:"industry"`
	Website        string   `json:"website"`
	Description    string   `json:"description"`
	CEO            string   `json:"CEO"`
	SecurityName   string   `json:"securityName"`
	IssueType      string   `json:"issueType"`
	Sector         string   `json:"sector"`
	PrimarySicCode int      `json:"primarySicCode"`
	Employees      int      `json:"employees"`
	Tags           []string `json:"tags"`
	Address        string   `json:"address"`
	Address2       string   `json:"address2"`
	State          string   `json:"state"`
	City           string   `json:"city"`
	Zip            string   `json:"zip"`
	Country        string   `json:"country"`
	Phone          string   `json:"phone"`
}

type CompanyKeyStats struct {
	CompanyName               string  `json:"companyName"`
	Marketcap                 int64   `json:"marketcap"`
	Week52High                float64 `json:"week52high"`
	Week52Low                 float64 `json:"week52low"`
	Week52HighSplitAdjustOnly float64 `json:"week52highSplitAdjustOnly"`
	Week52LowSplitAdjustOnly  float64 `json:"week52lowSplitAdjustOnly"`
	Week52Change              float64 `json:"week52change"`
	SharesOutstanding         int     `json:"sharesOutstanding"`
	Float                     int     `json:"float"`
	Avg10Volume               int     `json:"avg10Volume"`
	Avg30Volume               int     `json:"avg30Volume"`
	Day200MovingAvg           float64 `json:"day200MovingAvg"`
	Day50MovingAvg            float64 `json:"day50MovingAvg"`
	Employees                 int     `json:"employees"`
	TtmEPS                    float64 `json:"ttmEPS"`
	TtmDividendRate           float64 `json:"ttmDividendRate"`
	DividendYield             float64 `json:"dividendYield"`
	NextDividendDate          string  `json:"nextDividendDate"`
	ExDividendDate            string  `json:"exDividendDate"`
	NextEarningsDate          string  `json:"nextEarningsDate"`
	PeRatio                   float64 `json:"peRatio"`
	Beta                      float64 `json:"beta"`
	MaxChangePercent          float64 `json:"maxChangePercent"`
	Year5ChangePercent        float64 `json:"year5ChangePercent"`
	Year2ChangePercent        float64 `json:"year2ChangePercent"`
	Year1ChangePercent        float64 `json:"year1ChangePercent"`
	YtdChangePercent          float64 `json:"ytdChangePercent"`
	Month6ChangePercent       float64 `json:"month6ChangePercent"`
	Month3ChangePercent       float64 `json:"month3ChangePercent"`
	Month1ChangePercent       float64 `json:"month1ChangePercent"`
	Day30ChangePercent        float64 `json:"day30ChangePercent"`
	Day5ChangePercent         float64 `json:"day5ChangePercent"`
}

type ChartEntry struct {
	Date           string  `json:"date"`
	Close          float64 `json:"close"`
	Volume         int     `json:"volume"`
	Change         float64 `json:"change"`
	ChangePercent  float64 `json:"changePercent"`
	ChangeOverTime float64 `json:"changeOverTime"`
}
