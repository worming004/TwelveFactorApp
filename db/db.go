package db

import (
	"context"
	"os"
	"time"

	"github.com/worming004/TwelveFactorApp/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventRepository struct {
	client *mongo.Client
}

func NewDefaultMongoClient() (*mongo.Client, error) {
	mongoAddress := os.Getenv("TWELVE_DB_MONGO_ADDRESS")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoAddress))
	if err != nil {
		return nil, err
	}

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(timeout)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoEventRepository(client *mongo.Client) *MongoEventRepository {
	return &MongoEventRepository{client}
}

func (m *MongoEventRepository) CreateEvent(evt server.Event) error {
	collection := m.client.Database("event").Collection("event")
	bs := toBson(evt)
	_, err := collection.InsertOne(context.Background(), bs)
	return err
}

func toBson(evt server.Event) bson.M {
	return bson.M{"subject": evt.Subject}
}
