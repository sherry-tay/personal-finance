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
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"regularMarketPrice"`
	Currency string  `json:"currency"`
}

// CurrencyChart from Yahoo
type CurrencyChart struct {
	Chart CurrencyResult `json:"chart"`
}

// CurrencyResult in CurrencyChart
type CurrencyResult struct {
	Result []CurrencyMeta `json:"result"`
}

// CurrencyMeta in CurrencyResult
type CurrencyMeta struct {
	Meta Currency `json:"meta"`
}

// Currency in CurrencyMeta
type Currency struct {
	BaseCurrency string  `json:"currency"`
	Price        float64 `json:"regularMarketPrice"`
}
