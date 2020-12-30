package main

import (
	"fmt"

	"personal.finance/internal/firebase"
	"personal.finance/internal/sgx"
)

func main() {
	averagePrice := firebase.GetHoldings()
	for key, value := range averagePrice {
		current := sgx.GetCurrentPrice(key)
		fmt.Printf("Current %v, Average %v, Difference %v\n", current, value, current - value)
	}
}