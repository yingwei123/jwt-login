package mongodb

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mongoTest.io/server"
)

type ObjectID [12]byte

func (m MongoDBClient) CreateItem(item server.Item) (string, error) {
	itemCollection := m.client.Database(m.database).Collection("item")

	userResult, err := itemCollection.InsertOne(nil, item)
	if err != nil {
		return "", err
	}

	oid, err := insertOneID(userResult)
	if err != nil {
		return oid, err
	}

	return oid, nil
}

func insertOneID(fuc *mongo.InsertOneResult) (string, error) {
	if oid, ok := fuc.InsertedID.(primitive.ObjectID); ok {
		return oid.String(), nil
	}

	return "", errors.New("can't Get ObjectID of Item")
}
