package routes

import (
	"log" // log package
	"context" // context package
	"time" // time package
	"fmt" // fmt package

	"go.mongodb.org/mongo-driver/mongo" // driver package
	"go.mongodb.org/mongo-driver/mongo/options" // options package
)

func DBinstance() *mongo.Client { //creates a new mongo instance
	MongoDb := "mongodb://localhost:27017/caloriesdb"

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb)) //creates a new client
	if err!= nil {
        log.Fatal(err)
    }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //cancel the connection
	defer cancel()
	err = client.Connect(ctx) //connect to the database
	if err!= nil {
        log.Fatal(err)
    }
	fmt.Println("Connected to database")
	return client
}

var Client *mongo.Client = DBinstance() // create a new mongo instance

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection { // open a collection) *mongo.Collection { // open a collection)
	var collection *mongo.Collection = client.Database("caloriesdb").Collection(collectionName)
	return collection
}	