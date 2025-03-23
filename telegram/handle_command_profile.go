package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandProfile(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat

	user, err := bot.mc.UserGetByChatID(chat.Id)
	if err != nil {
		return err
	}

	if user.UserType == "driver" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"*Ваш профиль*"+"\n\n"+
				"Имя: "+user.Name+"\n"+
				"Роль: водитель",
			&gotgbot.SendMessageOpts{
				ParseMode: "markdown",
			},
		)
		return err
	}

	if user.UserType == "passenger" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"*Ваш профиль*"+"\n\n"+
				"Имя: "+user.Name+"\n"+
				"Роль: пассажир",
			&gotgbot.SendMessageOpts{
				ParseMode: "markdown",
			},
		)
		return err
	}

	return bot.onboard(&chat, user)
}
