package mongo

import (
	"context"

	"github.com/oybek/p24/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mc *MongoClient) GetGroupLastMessageId() (int64, error) {
	var state model.State
	err := mc.state.FindOne(context.Background(), bson.M{}).Decode(&state)
	return state.GroupLastMessageId, err
}

func (mc *MongoClient) SetGroupLastMessageId(id int64) error {
	_, err := mc.state.UpdateOne(
		context.Background(),
		bson.M{},
		bson.M{"$set": bson.M{"group_last_message_id": id}},
		options.Update().SetUpsert(true),
	)
	return err
}
