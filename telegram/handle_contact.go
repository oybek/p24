package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleContact(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat

	contact := ctx.EffectiveMessage.Contact

	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	user.Phone = contact.PhoneNumber
	err = bot.mc.UserUpdate(user)
	if err != nil {
		return err
	}

	return bot.onboard(chat, user)
}
