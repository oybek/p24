package telegram

import (
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const groupId = -1002626938267

func (bot *Bot) handleStartCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat

	user, err := bot.GetOrCreateUser(chat)
	if err != nil {
		return err
	}

	if user == nil || user.UserType == "" {
		_, err = bot.tg.SendMessage(
			chat.Id,
			fmt.Sprintf("Здравствуйте, %s!", chat.FirstName),
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{
					RemoveKeyboard: true,
				},
			},
		)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return bot.onboard(chat, user)
}
