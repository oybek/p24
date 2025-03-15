package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	client    *mongo.Client
	users     *mongo.Collection
	trips     *mongo.Collection
	cityNames *mongo.Collection
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
		client:    mongoClient,
		users:     mongoClient.Database("poputka").Collection("users"),
		trips:     mongoClient.Database("poputka").Collection("trips"),
		cityNames: mongoClient.Database("poputka").Collection("city_names"),
	}, nil
}

func (mc *MongoClient) CityNamesGet() (map[string]string, error) {
	var cityNames map[string]string
	err := mc.cityNames.FindOne(context.Background(), bson.M{}).Decode(&cityNames)
	if err != nil {
		return nil, err
	}
	return cityNames, nil
}
