package mongo

import (
	"context"
	"time"

	"github.com/oybek/p24/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (mc *MongoClient) TripCreate(
	trip *model.Trip,
) error {
	_, err := mc.trips.InsertOne(context.Background(), trip)
	return err
}

func (mc *MongoClient) TripFind(
	cityA string,
	cityB string,
	date time.Time,
) ([]model.Trip, error) {
	ctx := context.Background()

	date = date.Truncate(24 * time.Hour)
	filter := bson.M{
		"city_a": cityA,
		"city_b": cityB,
		"state":  "active",
		"start_time": bson.M{
			"$gte": date,
			"$lt":  date.Add(24 * time.Hour),
		},
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
