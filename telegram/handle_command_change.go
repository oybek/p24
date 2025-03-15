package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandChange(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	_, err := bot.tg.SendMessage(
		chat.Id,
		"Уточните, вы пассажир или водитель?",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbSelectRole(),
		},
	)
	return err
}
