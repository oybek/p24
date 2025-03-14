package mongo

import (
	"context"

	"github.com/oybek/p24/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mc *MongoClient) UserCreate(
	user *model.User,
) error {
	_, err := mc.users.InsertOne(context.Background(), user)
	return err
}

func (mc *MongoClient) UserUpdate(
	user *model.User,
) error {
	_, err := mc.users.UpdateOne(
		context.Background(),
		bson.M{"chat_id": user.ChatID},
		bson.M{"$set": user.BsonM()},
	)
	return err
}

func (mc *MongoClient) UserGetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*model.User, error) {
	var user model.User
	err := mc.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (mc *MongoClient) UserGetByChatID(
	chatID int64,
) (*model.User, error) {
	var user model.User
	err := mc.users.FindOne(
		context.Background(),
		bson.M{"chat_id": chatID},
	).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
