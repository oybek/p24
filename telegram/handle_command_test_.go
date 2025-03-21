package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandTest(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	_, err := bot.tg.SendMessage(
		chat.Id,
		"Ok",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{Text: "Поиск", WebApp: &gotgbot.WebAppInfo{Url: "https://oybek.github.io/p24-wa?user_type=search"}},
				}},
			},
		},
	)
	return err
}
