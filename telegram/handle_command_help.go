package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandHelp(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat

	user, err := bot.mc.UserGetByChatID(chat.Id)
	if err != nil {
		return err
	}

	return bot.onboard(&chat, user)
}
