package rest

import (
	"github.com/oybek/p24/mongo"
	"github.com/oybek/p24/telegram"
)

type Rest struct {
	bot *telegram.Bot
	mc  *mongo.MongoClient
}

func New(bot *telegram.Bot, mc *mongo.MongoClient) *Rest {
	return &Rest{
		bot: bot,
		mc:  mc,
	}
}
