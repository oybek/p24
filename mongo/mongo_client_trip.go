package mongo

import (
	"context"
	"log"

	"github.com/oybek/p24/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mc *MongoClient) TripCreate(
	trip *model.Trip,
) (primitive.ObjectID, error) {
	id, err := mc.trips.InsertOne(context.Background(), trip)
	return id.InsertedID.(primitive.ObjectID), err
}

func (mc *MongoClient) TripGetByID(
	tripID primitive.ObjectID,
) (*model.Trip, error) {
	var trip model.Trip
	err := mc.trips.FindOne(context.Background(), bson.M{"_id": tripID}).Decode(&trip)
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (mc *MongoClient) TripUpdateMessageID(
	tripID primitive.ObjectID,
	messageId int64,
) error {
	res, err := mc.trips.UpdateByID(
		context.Background(),
		tripID,
		bson.M{"$set": bson.M{"message_id": messageId}},
	)
	log.Printf("TripUpdateMessageID: %#v", res)
	return err
}

func (mc *MongoClient) TripDisable(
	tripID primitive.ObjectID,
) error {
	res, err := mc.trips.UpdateByID(
		context.Background(),
		tripID,
		bson.M{"$set": bson.M{"state": "disabled"}},
	)
	log.Printf("TripDisable: %#v", res)
	return err
}

func (mc *MongoClient) TripFind(
	userType string,
	cityA string,
	cityB string,
) ([]model.Trip, error) {
	ctx := context.Background()

	// TODO: exclude yesterday, and past trips
	filter := bson.M{
		"city_a":    cityA,
		"city_b":    cityB,
		"user_type": userType,
	}

	cursor, err := mc.trips.Find(ctx, filter)
	var trips []model.Trip
	for cursor.Next(ctx) {
		var trip model.Trip
		if err := cursor.Decode(&trip); err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	if err != nil {
		return nil, err
	}

	return trips, nil
}
