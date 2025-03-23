package telegram

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandPassenger(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	user.UserType = "passenger"

	err = bot.mc.UserUpdate(user)
	if err != nil {
		return err
	}

	cb := ctx.CallbackQuery
	cb.Message.EditText(b,
		"Вы теперь пассажир",
		&gotgbot.EditMessageTextOpts{},
	)
	cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	time.Sleep(1 * time.Second)

	return bot.onboard(chat, user)
}
