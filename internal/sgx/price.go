package sgx

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"os"
)

const categoriesJSONFilePath = "data/category.json"

// GetCurrentPrice gets the current price of the stock from SGX
func GetCurrentPrice(code string) (float64, error) {
	jsonFile, err := os.Open(categoriesJSONFilePath)
	if err != nil {
		fmt.Printf("Failed to open category json file: %v", err)
		return 0.0, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var categories map[string]string
	json.Unmarshal([]byte(byteValue), &categories)
	if _, ok := categories[code]; !ok {
		fmt.Printf("Failed to find category info for code %v: %v", code, err)
		return 0.0, err
	}

	resp, err := http.Get("https://api.sgx.com/securities/v1.1/" + categories[code] + "/code/" + code)
	if err != nil {
		fmt.Printf("Failed to get securities data: %v", err)
		return 0.0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read securities data response body: %v", err)
		return 0.0, err
	}

	var message SecuritiesData
	if err := json.Unmarshal(body, &message); err != nil {
		fmt.Printf("Failed to unmarshal securities data: %v", err)
		return 0.0, err
	}
	defer resp.Body.Close()
	return message.Data.Prices[0].LastTraded, nil
}
