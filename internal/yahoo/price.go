package yahoo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	endpoint = "https://query1.finance.yahoo.com/v8/finance/chart/"
)

// GetCurrentPrice gets the current price of the stock from Yahoo
func GetCurrentPrice(code string) (float64, error) {
	resp, err := http.Get(endpoint + code)
	if err != nil {
		fmt.Printf("Failed to get securities data: %v", err)
		return 0.0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Failed to read securities data response body: %v", err)
		return 0.0, err
	}

	var message SecuritiesChart
	if err := json.Unmarshal(body, &message); err != nil {
		fmt.Printf("Failed to unmarshal securities data: %v", err)
		return 0.0, err
	}

	if len(message.Chart.Result) < 1 {
		return 0.0, fmt.Errorf("No price value found for %v", code)
	}
	return message.Chart.Result[0].Meta.Price, nil
}
