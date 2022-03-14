package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client   *mongo.Client
	database string
	context  context.Context
	cancel   context.CancelFunc
}

func CreateMongoClient(atlasURI string) (MongoDBClient, error) {
	clientOptions := options.Client().ApplyURI(atlasURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return MongoDBClient{}, err
	}

	mongoDbClient := MongoDBClient{client: client, database: "myFirstDatabase", context: ctx, cancel: cancel}

	mongoDbClient.CreateIndex("user", "email", true)

	return mongoDbClient, nil
}

func (m MongoDBClient) Disconnect() {
	m.client.Disconnect(m.context)
	defer m.cancel()
	log.Println("Disconnecting Mongodb")
}
