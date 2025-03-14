package telegram

import "github.com/PaulSonOfLars/gotgbot/v2/ext"

func handleCommandHelp(bot *Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat

	user, err := bot.mc.UserGetByChatID(chat.Id)
	if err != nil {
		return err
	}

	return bot.onboard(&chat, user)
}
