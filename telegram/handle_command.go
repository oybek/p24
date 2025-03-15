package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/samber/lo"
)

var commands = []gotgbot.BotCommand{
	{Command: "change", Description: "Сменить роль"},
	{Command: "profile", Description: "Мой профиль"},
	{Command: "help", Description: "Помощь"},
}

func (lp *Bot) SetupCommands() error {
	botCommands := lo.Map(commands, func(cmd gotgbot.BotCommand, _ int) gotgbot.BotCommand {
		return gotgbot.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		}
	})

	_, err := lp.tg.SetMyCommands(botCommands, nil)
	return err
}
