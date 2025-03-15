package telegram

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandGroupUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Printf("handleCommandGroupUpdate")
	_, err := bot.tg.SendMessage(
		groupId,
		"Создайте объявление через бота",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbOpenBot(),
		},
	)
	return err
}
