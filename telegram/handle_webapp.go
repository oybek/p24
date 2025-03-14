package telegram

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (lp *Bot) handleWebAppData(b *gotgbot.Bot, ctx *ext.Context) error {
	webAppData := ctx.EffectiveMessage.WebAppData
	log.Printf("Got data from webapp: %s", webAppData.Data)

	return nil
}
