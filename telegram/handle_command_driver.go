package telegram

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const UserTypeDriver = "driver"

func (bot *Bot) handleCommandDriver(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	user.UserType = UserTypeDriver
	user.Phone = ""
	user.CarPhoto = ""

	err = bot.mc.UserUpdate(user)
	if err != nil {
		return err
	}

	ctx.CallbackQuery.Message.EditText(b,
		"Вы теперь водитель",
		&gotgbot.EditMessageTextOpts{},
	)
	time.Sleep(1 * time.Second)

	return bot.onboardDriver(user)
}
