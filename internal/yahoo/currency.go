package yahoo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

// GetCurrency gets the current price of the buyingCurrency in units of sellingCurrency from Yahoo
func GetCurrency(buyingCurrency, sellingCurrency string) (float64, error){
	endpoint := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s%s=X?interval=1d", buyingCurrency, sellingCurrency)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request to get currency data: %v", err)
		return 0.0, err
	}
	req.Header.Set("User-Agent", uuid.NewString())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to get currency data: %v", err)
		return 0.0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Failed to read currency data response body: %v", err)
		return 0.0, err
	}

	var message CurrencyChart
	if err := json.Unmarshal(body, &message); err != nil {
		fmt.Printf("Failed to unmarshal currency data: %v", err)
		return 0.0, err
	}

	if len(message.Chart.Result) < 1 {
		return 0.0, fmt.Errorf("No price value found for %v in terms of %v", buyingCurrency, sellingCurrency)
	}
	return message.Chart.Result[0].Meta.Price, nil
}
