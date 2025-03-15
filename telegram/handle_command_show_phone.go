package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandShowPhone(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.CallbackQuery
	cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "0559171775",
		ShowAlert: true,
	})
	return nil
}
