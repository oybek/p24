package rest

import "github.com/oybek/p24/mongo"

type Rest struct {
	mc *mongo.MongoClient
}

func New(mongoClient *mongo.MongoClient) *Rest {
	return &Rest{
		mc: mongoClient,
	}
}
