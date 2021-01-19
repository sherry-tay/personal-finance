package firestore

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	serviceAccountFilePath = "personal-finance-admin.json"
	stocksCollectionName   = "stocks"
)

func getFirebaseClient() (*firestore.Client, context.Context, error) {
	ctx := context.Background()

	var app *firebase.App
	var err error

	if os.Getenv("IS_GCP") != "" {
		conf := &firebase.Config{ProjectID: os.Getenv("PROJECT_ID")}
		app, err = firebase.NewApp(ctx, conf)
	} else {
		sa := option.WithCredentialsFile(serviceAccountFilePath)
		app, err = firebase.NewApp(ctx, nil, sa)
	}
	
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	fmt.Println("Successfully obtained Firestore client")
	return client, ctx, nil
}

// GetHoldings stored in Firestore
func GetHoldings() (map[string]Stock, error) {
	client, ctx, err := getFirebaseClient()
	if err != nil {
		return nil, err
	}

	iter := client.Collection(stocksCollectionName).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		fmt.Printf("Failed to get documents: %v", err)
		return nil, err
	}
	fmt.Println("Successfully obtained stock documents")

	var stocksList []Stock
	for _, doc := range docs {
		var s Stock
		if err := doc.DataTo(&s); err != nil {
			fmt.Printf("Failed to transform: %v", err)
			return nil, err
		}
		stocksList = append(stocksList, s)
	}

	fmt.Printf("Stocks in holdings: %v", stocksList)

	defer client.Close()

	return getAveragePrice(stocksList), nil
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
		averaged[code] = Stock{
			Code:   code,
			Price:  sumProduct / float64(volume),
			Volume: volume,
		}
	}

	fmt.Printf("Stocks in holdings (averaged): %v", averaged)

	return averaged
}

// AddHoldings to Firestore
func AddHoldings(id string, stock Stock) error {
	client, ctx, err := getFirebaseClient()
	if err != nil {
		return err
	}

	if writeResult, err := client.Collection(stocksCollectionName).Doc(id).Set(ctx, stock); err != nil {
		fmt.Printf("Failed to set documents: %v", err)
		return err
	} else {
		fmt.Printf("Successfully added holdings at %v", writeResult)
		return nil
	}
}

// Stock as recorded in Firestore
type Stock struct {
	Code     string    `firestore:"code"`
	Price    float64   `firestore:"price"`
	Volume   int       `firestore:"volume"`
	Date     time.Time `firestore:"date"`
	StoredIn string    `firestore:"in"`
}
