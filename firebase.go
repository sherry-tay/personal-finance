package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var serviceAccountFilePath = "personal-finance-admin.json"
var stocksCollectionName = "stocks"

func main() {
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

	defer client.Close()
}

type Stock struct {
	Code    string  	`firestore:"Code"`
	Price 	float64 	`firestore:"Price"`
	Volume 	int64 		`firestore:"Volume"`
}