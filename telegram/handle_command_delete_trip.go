package telegram

import (
	"errors"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (bot *Bot) handleCommandDeleteTrip(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.CallbackQuery
	hex, found := strings.CutPrefix(cb.Data, "/del")
	if !found {
		return errors.New("/del command handle error")
	}

	tripID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return err
	}

	// delete in user chat
	cb.Message.Delete(b, &gotgbot.DeleteMessageOpts{})

	trip, err := bot.mc.TripGetByID(tripID)
	if err != nil {
		return err
	}

	// delete in group
	b.DeleteMessage(groupId, trip.MessageId, &gotgbot.DeleteMessageOpts{})

	return bot.mc.TripDisable(tripID)
}
