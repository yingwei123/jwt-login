package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"mongoTest.io/server"
)

func (m MongoDBClient) CreateNewUser(user server.NewUser) (string, error) {
	userCollection := m.client.Database(m.database).Collection("user")

	userResult, err := userCollection.InsertOne(nil, user)
	if err != nil {
		return "", err
	}

	oid, err := insertOneID(userResult)
	if err != nil {
		return oid, err
	}

	return oid, nil
}

func (m MongoDBClient) FindUserByEmail(email string) (server.NewUser, error) {
	user := server.NewUser{}
	filter := bson.M{"email": bson.M{"$eq": email}}
	err := m.client.Database(m.database).Collection("user").FindOne(nil, filter).Decode(&user)
	if err != nil {
		println(err.Error())
	}

	println(fmt.Sprintf("%v", user))

	return user, nil
}

//create unique collection name
func (m MongoDBClient) CreateIndex(collectionName string, field string, unique bool) bool {

	// 1. Lets define the keys for the index we want to create
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	collection := m.client.Database(m.database).Collection(collectionName)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}
