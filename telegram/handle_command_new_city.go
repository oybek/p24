package telegram

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) handleCommandNewCity(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	if !slices.Contains(agentIds, chat.Id) {
		return nil
	}

	text := ctx.EffectiveMessage.Text
	data, found := strings.CutPrefix(text, "/new_city")
	if !found {
		return errors.New("/new_city command handle error 0")
	}

	datas := strings.Split(data, "/")
	if len(datas) != 2 {
		return errors.New("/new_city command handle error 1")
	}
	key := strings.TrimSpace(datas[0])
	value := strings.TrimSpace(datas[1])

	bot.cityNames.Set(key, value)
	err := bot.mc.CityNamesAdd(key, value)
	if err != nil {
		return err
	}

	_, err = bot.tg.SendMessage(
		chat.Id,
		fmt.Sprintf("Добавлен город %s: %s", key, value),
		&gotgbot.SendMessageOpts{},
	)
	return err
}
