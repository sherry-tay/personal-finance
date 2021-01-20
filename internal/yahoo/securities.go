package yahoo

// SecuritiesChart from Yahoo
type SecuritiesChart struct {
	Chart SecuritiesResult `json:"chart"`
}

// SecuritiesResult in SecuritiesChart
type SecuritiesResult struct {
	Result []SecuritiesMeta `json:"result"`
}

// SecuritiesMeta in SecuritiesResult
type SecuritiesMeta struct {
	Meta SecuritiesPrice `json:"meta"`
}

// SecuritiesPrice in SecuritiesMeta
type SecuritiesPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"regularMarketPrice"`
}
