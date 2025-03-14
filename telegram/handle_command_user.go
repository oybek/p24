package telegram

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const UserTypeUser = "user"

func (bot *Bot) handleCommandUser(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	user.UserType = UserTypeUser

	err = bot.mc.UserUpdate(user)
	if err != nil {
		return err
	}

	ctx.CallbackQuery.Message.EditText(b,
		"Вы теперь пассажир",
		&gotgbot.EditMessageTextOpts{},
	)
	time.Sleep(1 * time.Second)

	return bot.onboard(chat, user)
}
