package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/p24/model"
)

func (bot *Bot) handleWebAppTrip(chat *gotgbot.Chat, trip *model.Trip) error {
	trip.ChatID = chat.Id
	trip.State = "active"

	err := bot.mc.TripCreate(trip)
	if err != nil {
		return err
	}

	_, err = bot.tg.SendMessage(
		chat.Id,
		"*Поездка создана*"+"\n\n"+
			Show(trip),
		&gotgbot.SendMessageOpts{
			ParseMode: "markdown",
		},
	)

	return err
}
