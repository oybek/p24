package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    int64              `bson:"chat_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	UserType  string             `bson:"user_type,omitempty"`
	Phone     string             `bson:"phone,omitempty"`
	CarPhoto  string             `bson:"car_photo,omitempty"`
	StartTime time.Time          `bson:"start_time,omitempty"`
}

func (user *User) BsonM() bson.M {
	bsonM := bson.M{}
	bsonM["chat_id"] = user.ChatID
	bsonM["name"] = user.Name
	bsonM["user_type"] = user.UserType
	bsonM["phone"] = user.Phone
	bsonM["car_photo"] = user.CarPhoto
	bsonM["start_time"] = user.StartTime
	return bsonM
}
