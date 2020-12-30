package firebase

import (
	"context"
	"fmt"
	"log"
	
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var serviceAccountFilePath = "../internal/firebase/personal-finance-admin.json"
var stocksCollectionName = "stocks"

func Initialise() {
	ctx := context.Background()
	sa := option.WithCredentialsFile(serviceAccountFilePath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil { 
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	iter := client.Collection(stocksCollectionName).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		log.Fatalf("Failed to get documents: %v", err)
	}

	var stocksList []Stock
	for _, doc := range docs {
		var s Stock
		if err := doc.DataTo(&s); err != nil {
			log.Fatalf("Failed to transform: %v", err)
		}
		stocksList = append(stocksList, s)
	}

	fmt.Println(stocksList)

	fmt.Println(getAveragePrice(stocksList))

	defer client.Close()
}

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

type Stock struct {
	Code    string  	`firestore:"Code"`
	Price 	float64 	`firestore:"Price"`
	Volume 	int 		`firestore:"Volume"`
}