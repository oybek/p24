package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handlePhoto(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	photo := ctx.EffectiveMessage.Photo[0].FileId

	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	user.CarPhoto = photo
	err = bot.mc.UserUpdate(user)
	if err != nil {
		return err
	}

	return bot.onboard(chat, user)
}
