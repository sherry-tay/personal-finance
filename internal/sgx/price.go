package sgx

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var categoriesJsonFilePath = "../data/category.json"

func GetCurrentPrice(code string) float64 {
	jsonFile, err := os.Open(categoriesJsonFilePath)
	if err != nil {
		log.Fatalf("Failed to open category json file: %v", err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var categories map[string]string
	json.Unmarshal([]byte(byteValue), &categories)
	if _, ok := categories[code]; !ok {
		log.Fatalf("Failed to find category info for code %v: %v", code, err)
	}

	resp, err := http.Get("https://api.sgx.com/securities/v1.1/" + categories[code] + "/code/" + code)
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