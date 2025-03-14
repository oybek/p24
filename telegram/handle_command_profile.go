package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func handleCommandProfile(bot *Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat

	user, err := bot.mc.UserGetByChatID(chat.Id)
	if err != nil {
		return err
	}

	if user.UserType == "driver" {
		_, err := bot.tg.SendPhoto(
			user.ChatID,
			gotgbot.InputFileByID(user.CarPhoto),
			&gotgbot.SendPhotoOpts{
				Caption: "*Ваш профиль*" + "\n\n" +
					"Имя: " + user.Name + "\n" +
					"Номер телефона: " + user.Phone,
				ParseMode: "markdown",
			},
		)
		return err
	}

	if user.UserType == "user" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"*Ваш профиль*"+"\n\n"+
				"Имя: "+user.Name+"\n",
			&gotgbot.SendMessageOpts{
				ParseMode: "markdown",
			},
		)
		return err
	}

	return bot.onboard(&chat, user)
}
