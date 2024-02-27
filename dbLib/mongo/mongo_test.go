package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestDialMongoDB(t *testing.T) {
	mongoClientPtr := Dial("mongodb://localhost:27017")
	defer mongoClientPtr.Close()

	dbPtr := mongoClientPtr.GetDatabase("test")
	collectionPtr := dbPtr.Collection("test")

	// insert
	insertResult, err := collectionPtr.InsertOne(context.Background(), bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		t.Errorf("TestDialMongoDB InsertOne err:%s", err.Error())
	}
	t.Logf("TestDialMongoDB InsertOne insertResult:%+v", insertResult)

	// find
	var result bson.M
	err = collectionPtr.FindOne(context.Background(), bson.D{{"name", "pi"}}).Decode(&result)
	if err != nil {
		t.Errorf("TestDialMongoDB FindOne err:%s", err.Error())
	}
	t.Logf("TestDialMongoDB FindOne result:%+v", result)

	// update
	updateResult, err := collectionPtr.UpdateOne(context.Background(), bson.D{{"name", "pi"}}, bson.D{{"$set", bson.D{{"value", 3.1415926}}}})
	if err != nil {
		t.Errorf("TestDialMongoDB UpdateOne err:%s", err.Error())
	}
	t.Logf("TestDialMongoDB UpdateOne updateResult:%+v", updateResult)

	// delete
	deleteResult, err := collectionPtr.DeleteOne(context.Background(), bson.D{{"name", "pi"}})
	if err != nil {
		t.Errorf("TestDialMongoDB DeleteOne err:%s", err.Error())
	}
	t.Logf("TestDialMongoDB DeleteOne deleteResult:%+v", deleteResult)
}
