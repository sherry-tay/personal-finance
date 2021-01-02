package main

import (
	"fmt"
	"math"

	"personal.finance/internal/firebase"
	"personal.finance/internal/sgx"
)

func main() {
	totalPortfolio, totalInvested := 0.0, 0.0
	averagePrice := firebase.GetHoldings()
	for key, value := range averagePrice {
		current := sgx.GetCurrentPrice(key)
		diff, percentage, profit := getProfit(current, value.Price, value.Volume)
		fmt.Printf("%v: Current %v, Average %v, Difference %v (%v%%), Volume %v, Profit %v\n", key, current, value.Price, diff, percentage, value.Volume, profit)
		totalPortfolio += current * float64(value.Volume)
		totalInvested += value.Price * float64(value.Volume)
	}
	totalDiff, totalPercentage, _ := getProfit(totalPortfolio, totalInvested, 0)
	fmt.Printf("Current portfolio: %v\n", round(totalPortfolio))
	fmt.Printf("Total invested: %v\n", round(totalInvested))
	fmt.Printf("Capital gains: %v (%v%%)\n", totalDiff, totalPercentage)
}

func getProfit(current, average float64, volume int) (float64, float64, float64) {
	diff := (current - average)
	return round(diff), round(diff/average*100), round(diff * float64(volume))
}

func round(val float64) float64 {
	precisionFactor := 1000.0
	return math.Round(val * precisionFactor)/precisionFactor
}