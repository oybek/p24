package telegram

import (
	"bytes"
	"slices"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
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
		"Уточните, вы пассажир или водитель?",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbSelectRole(),
		},
	)
	return err
}

func (bot *Bot) onboardDriver(user *model.User) error {
	if user.Phone == "" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"Чтобы стать проверенным водителем - поделитесь своим контактом",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: kbSendContact(),
			},
		)
		return err
	}

	if user.CarPhoto == "" {
		_, err := bot.tg.SendMessage(
			user.ChatID,
			"Теперь отправьте фото своей машины",
			&gotgbot.SendMessageOpts{},
		)
		return err
	}

	_, err := bot.tg.SendMessage(
		user.ChatID,
		"✅ Вы наш проверенный водитель!\n"+
			"Переходите в группу и находите попутчиков",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbOpenGroup(),
		},
	)
	return err
}

func (bot *Bot) onboardUser(user *model.User) error {
	_, err := bot.tg.SendMessage(
		user.ChatID,
		"Нажмите кнопку 'Создать карточку'",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbCreateTrip(slices.Contains(agentIds, user.ChatID)),
		},
	)
	return err
}

func (bot *Bot) deleteMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	_, err := b.DeleteMessage(chat.Id, ctx.EffectiveMessage.MessageId, &gotgbot.DeleteMessageOpts{})
	return err
}

func (bot *Bot) publishCard(
	chat *gotgbot.Chat,
	trip *model.Trip,
	card []byte,
) (*gotgbot.Message, error) {
	bot.tg.DeleteMessage(groupId, groupLastMessageId.Load(), &gotgbot.DeleteMessageOpts{})
	cardMessage, _ := bot.tg.SendPhoto(
		groupId,
		gotgbot.InputFileByReader("img.jpg", bytes.NewReader(card)),
		&gotgbot.SendPhotoOpts{
			ReplyMarkup: kbUnderCardInGroup(chat, trip),
		},
	)
	groupLastMessage, _ := bot.tg.SendMessage(
		groupId,
		"Создайте объявление через бота",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbOpenBot(),
		},
	)
	groupLastMessageId.Store(groupLastMessage.MessageId)
	return cardMessage, nil
}
