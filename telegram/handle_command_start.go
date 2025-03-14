package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func handleStartCommand(bot *Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat

	_, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	_, err = bot.tg.SendMessage(
		chat.Id,
		"Здравствуйте!",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{
				RemoveKeyboard: true,
			},
		},
	)
	if err != nil {
		return err
	}

	_, err = bot.tg.SendMessage(
		chat.Id,
		"Вы пассажир или водитель?",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbSelectRole(),
		},
	)

	return err
}
