package sgx

// SecuritiesData from SGX
type SecuritiesData struct {
	Data SecuritiesPrices `json:"data"`
}

// SecuritiesPrices in SecuritiesData
type SecuritiesPrices struct {
	Prices []SecuritiesPrice `json:"prices"`
}

// SecuritiesPrice in SecuritiesPrices
type SecuritiesPrice struct {
	Code       string  `json:"nc"`
	Type       string  `json:"type"`
	LastTraded float64 `json:"lt"`
}
