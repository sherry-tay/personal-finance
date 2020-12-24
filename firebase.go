package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/api/iterator"
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
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}

	defer client.Close()
}