package firebase

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var serviceAccountFilePath = "../internal/firebase/personal-finance-admin.json"
var stocksCollectionName = "stocks"

func getFirebaseClient() (*firestore.Client, context.Context) {
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

	return client, ctx
}

// GetHoldings stored in Firestore
func GetHoldings() map[string]Stock {
	client, ctx := getFirebaseClient()

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

	defer client.Close()

	return getAveragePrice(stocksList)
}

func getAveragePrice(list []Stock) map[string]Stock {
	collatedMap := make(map[string][]Stock)

	for _, stock := range list {
		collatedMap[stock.Code] = append(collatedMap[stock.Code], stock)
	}

	averaged := make(map[string]Stock)

	for code, stocks := range collatedMap {
		sumProduct := 0.0
		volume := 0
		for _, stock := range stocks {
			 sumProduct += stock.Price * float64(stock.Volume)
			 volume += stock.Volume
		}
		averaged[code] = Stock {
			Code: code,
			Price: sumProduct / float64(volume),
			Volume: volume,
		}
	}

	return averaged
}

// AddHoldings to Firestore
func AddHoldings(id string, stock Stock) {
	client, ctx := getFirebaseClient()

	_, err := client.Collection(stocksCollectionName).Doc(id).Set(ctx, stock)
	if err != nil {
		log.Fatalf("Failed to set documents: %v", err)
	}
}

// Stock as recorded in Firestore
type Stock struct {
	Code    	string  	`firestore:"code"`
	Price 		float64 	`firestore:"price"`
	Volume 		int 		`firestore:"volume"`
	Date		time.Time	`firestore:"date"`
	StoredIn	string		`firestore:"in"`
}