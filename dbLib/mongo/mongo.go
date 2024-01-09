package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	client *mongo.Client
}

func Dial(connectUri string) *MongoClient {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectUri).SetServerAPIOptions(serverApi)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoClientPtr, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	clientPtr := &MongoClient{client: mongoClientPtr}

	return clientPtr
}

func (this *MongoClient) GetClient() *mongo.Client {
	return this.client
}

func (this *MongoClient) Close() {
	this.client.Disconnect(context.Background())
}

func (this *MongoClient) GetDatabase(databaseName string) *mongo.Database {
	return this.client.Database(databaseName)
}

func (this *MongoClient) GetCollection(databaseName string, collectionName string) *mongo.Collection {
	return this.GetDatabase(databaseName).Collection(collectionName)
}
