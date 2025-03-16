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

const database = "p24"

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
		users:     mongoClient.Database(database).Collection("users"),
		trips:     mongoClient.Database(database).Collection("trips"),
		cityNames: mongoClient.Database(database).Collection("city_names"),
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

func (mc *MongoClient) CityNamesAdd(key string, value string) error {
	_, err := mc.cityNames.UpdateOne(
		context.Background(),
		bson.M{},
		bson.M{"$set": bson.M{key: value}},
	)
	return err
}
