package telegram

import (
	"errors"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/oybek/p24/model"
)

func (bot *Bot) handleWebAppData(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	webAppData := ctx.EffectiveMessage.WebAppData
	if webAppData == nil {
		return errors.New("empty webapp data")
	}

	b.DeleteMessage(chat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{})
	log.Printf("Got data from webapp: %s", webAppData.Data)

	if trip, err := model.ParseAndValidate[model.Trip](webAppData.Data); err == nil {
		return bot.handleWebAppTrip(chat, trip)
	}

	return errors.New("invalid webapp data")
}
