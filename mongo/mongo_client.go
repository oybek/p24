package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	client   *mongo.Client
	database *mongo.Database
	users    *mongo.Collection
}

func NewMongoClient(ctx context.Context, URL string) (*MongoClient, error) {
	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(URL),
	)

	if err != nil {
		log.Fatalf("connection error :%v", err)
		return nil, err
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ping mongodb error :%v", err)
		return nil, err
	}

	return &MongoClient{
		client:   mongoClient,
		database: mongoClient.Database("poputka"),
		users:    mongoClient.Database("poputka").Collection("users"),
	}, nil
}
