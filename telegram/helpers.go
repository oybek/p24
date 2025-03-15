package telegram

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/p24/model"
)

func (bot *Bot) GetOrCreateUser(chat *gotgbot.Chat) (*model.User, error) {
	user, err := bot.mc.UserGetByChatID(chat.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &model.User{
			ChatID:    chat.Id,
			Name:      chat.FirstName,
			StartTime: time.Now(),
		}

		err = bot.mc.UserCreate(user)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func (bot *Bot) onboard(chat *gotgbot.Chat, user *model.User) error {
	if user.UserType == "driver" {
		return bot.onboardDriver(user)
	}

	if user.UserType == "user" {
		return bot.onboardUser(user)
	}

	_, err := bot.tg.SendMessage(
		chat.Id,
		"Вы пассажир или водитель?",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbShowPhone(),
		},
	)
	return err
}

func (bot *Bot) onboardDriver(user *model.User) error {
	if user.Phone == "" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"Поделитесь своим контактом",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: kbSendContact(),
			},
		)
		return err
	}

	if user.CarPhoto == "" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"Отправьте фото свой машины",
			&gotgbot.SendMessageOpts{},
		)
		return err
	}

	_, err := bot.tg.SendMessage(
		user.ChatID,
		"Нажмите кнопку 'Найти попутчиков'",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbFindPassengers(),
		},
	)
	return err
}

func (bot *Bot) onboardUser(user *model.User) error {
	_, err := bot.tg.SendMessage(
		user.ChatID,
		"Нажмите кнопку 'Создать поездку'",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbCreateTrip(),
		},
	)
	return err
}
