package sgx

type SecuritiesData struct {
	Data SecuritiesPrices `json:"data"`
} 

type SecuritiesPrices struct {
	Prices []SecuritiesPrice `json:"prices"`
}

type SecuritiesPrice struct {
	Code 		string 	`json:"nc"`
	Type 		string 	`json:"type"`
	LastTraded 	float64 `json:"lt"`
}