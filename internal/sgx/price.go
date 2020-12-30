package sgx

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetCurrentPrice(code, category string) float64 {
	resp, err := http.Get("https://api.sgx.com/securities/v1.1/" + category + "/code/" + code)
	if err != nil {
		log.Fatalf("Failed to get securities data: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var message SecuritiesData
	if err := json.Unmarshal(body, &message); err != nil {
		log.Fatalf("Failed to unmarshal securities data: %v", err)
	}
	defer resp.Body.Close()
	return message.Data.Prices[0].LastTraded
}