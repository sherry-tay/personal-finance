package main

import (
	"fmt"
	"math"

	"personal.finance/internal/firebase"
	"personal.finance/internal/sgx"
)

func main() {
	totalProfit := 0.0
	averagePrice := firebase.GetHoldings()
	for key, value := range averagePrice {
		current := sgx.GetCurrentPrice(key)
		diff, profit := getProfit(current, value.Price, value.Volume)
		fmt.Printf("%v: Current %v, Average %v, Difference %v, Volume %v, Profit %v\n", key, current, value.Price, diff, value.Volume, profit)
		totalProfit += profit
	}
	fmt.Printf("Capital gains: %v\n", totalProfit)
}

func getProfit(current, average float64, volume int) (float64, float64) {
	precisionFactor := 1000.0
	diff := (current - average) * precisionFactor
	return math.Round(diff)/precisionFactor, math.Round(diff * float64(volume))/precisionFactor
}