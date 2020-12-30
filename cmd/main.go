package main

import (
	"fmt"

	"personal.finance/internal/firebase"
	"personal.finance/internal/sgx"
)

func main() {
	firebase.Initialise()
	fmt.Println(sgx.GetCurrentPrice("G3B", "etfs"))
}