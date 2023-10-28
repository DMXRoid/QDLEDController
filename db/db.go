package db

import (
	"context"
	"fmt"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func Init() {
	var err error

	dsn := "mongodb://0.0.0.0:27017"

	opts := options.Client().ApplyURI(dsn)

	mongoClient, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		fmt.Println(fmt.Sprintf(":::MONGO ERROR::: %s", err))
	}

}

func GetClient() *mongo.Client {
	return mongoClient
}

func GetCollection(n string) *mongo.Collection {
	return mongoClient.Database("qdled").Collection(n)
}
