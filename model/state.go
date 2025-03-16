package model

type State struct {
	GroupLastMessageId int64 `bson:"group_last_message_id"`
}
