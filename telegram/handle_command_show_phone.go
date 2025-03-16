package telegram

import (
	"errors"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (bot *Bot) handleCommandShowPhone(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.CallbackQuery
	hex, found := strings.CutPrefix(cb.Data, "/show_phone")
	if !found {
		return errors.New("/show_phone command handle error")
	}

	user, err := bot.mc.UserGetByChatID(cb.From.Id)
	if err != nil {
		return err
	}
	if user.CarPhoto == "" || user.Phone == "" {
		_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Только проверенные водители могут видеть номера. Пройдите простую проверку через нашего бота, это займет всего пару минут! 😊",
			ShowAlert: true,
		})
		return err
	}

	tripID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return err
	}

	trip, err := bot.mc.TripGetByID(tripID)
	if err != nil {
		return err
	}

	_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text:      trip.Phone,
		ShowAlert: true,
	})
	return err
}
