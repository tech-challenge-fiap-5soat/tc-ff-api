package config

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func InitMongoDbConfiguration(ctx context.Context) (mongo.Client, error) {

	mongoCfg := GetMongoCfg()
	//uri := fmt.Sprintf("mongodb://%s:%s", mongoCfg.Host, mongoCfg.Port)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/?appName=%s&retryWrites=true&w=majority", mongoCfg.User, mongoCfg.Pass, mongoCfg.Host, mongoCfg.Database)

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client
	return *client, err
}

func GetMongoClient() (mongo.Client, error) {
	if mongoClient != nil {
		return *mongoClient, nil
	}
	return mongo.Client{}, errors.New("mongo client not initialized")
}
