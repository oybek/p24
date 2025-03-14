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
	Handler     func(*Bot, *ext.Context) error
}

var commands = []BotCommand{
	{Command: "start", Description: "Начать заново", Handler: handleStartCommand},
	{Command: "profile", Description: "Мой профиль", Handler: handleCommandProfile},
	{Command: "help", Description: "Помощь", Handler: handleCommandHelp},
}

func (lp *Bot) SetupCommands() error {
	botCommands := lo.Map(commands, func(cmd BotCommand, _ int) gotgbot.BotCommand {
		return gotgbot.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		}
	})

	_, err := lp.tg.SetMyCommands(botCommands, nil)
	return err
}

func (lp *Bot) handleCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	text := ctx.EffectiveMessage.Text

	for _, cmd := range commands {
		if strings.HasSuffix(text, cmd.Command) {
			return cmd.Handler(lp, ctx)
		}
	}

	return handleCommandHelp(lp, ctx)
}
