package main

func getAveragePrice(list []Stock) map[string]float64 {
	collatedMap := make(map[string][]Stock)

	for _, stock := range list {
		collatedMap[stock.Code] = append(collatedMap[stock.Code], stock)
	}

	averagedMap := make(map[string]float64)

	for code, stocks := range collatedMap {
		sumProduct := 0.0
		volume := 0
		for _, stock := range stocks {
			 sumProduct += stock.Price * float64(stock.Volume)
			 volume += stock.Volume
		}
		averagedMap[code] = sumProduct / float64(volume)
	}

	return averagedMap
}