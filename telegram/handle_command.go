package telegram

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/samber/lo"
)

type BotCommand struct {
	Command     string
	Description string
	Handler     func(*LongPoll, *ext.Context) error
}

var commands = []BotCommand{
	{Command: "create_apteka", Description: "Создать аптеку", Handler: createAptekaHandler},
	{Command: "help", Description: "Помощь", Handler: helpHandler},
}

func (lp *LongPoll) SetupCommands() error {
	botCommands := lo.Map(commands, func(cmd BotCommand, _ int) gotgbot.BotCommand {
		return gotgbot.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		}
	})

	_, err := lp.bot.SetMyCommands(botCommands, nil)
	return err
}

func (lp *LongPoll) handleCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	text := ctx.EffectiveMessage.Text

	for _, cmd := range commands {
		if strings.HasSuffix(text, cmd.Command) {
			return cmd.Handler(lp, ctx)
		}
	}

	return helpHandler(lp, ctx)
}

func createAptekaHandler(lp *LongPoll, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	createAptekaKeyboard := &gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: "Создать аптеку", WebApp: &gotgbot.WebAppInfo{Url: createAptekaWebAppUrl}},
			},
		},
	}
	_, err := lp.bot.SendMessage(chat.Id, TextCreateApteka,
		&gotgbot.SendMessageOpts{ReplyMarkup: createAptekaKeyboard})
	return err
}

func helpHandler(lp *LongPoll, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	return lp.sendText(chat.Id, TextDefault)
}
