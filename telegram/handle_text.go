package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	return bot.onboard(chat, user)
}
