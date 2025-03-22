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
		return bot.onboardUser(user)
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

func (bot *Bot) onboardUser(user *model.User) error {
	_, err := bot.tg.SendMessage(
		user.ChatID,
		"Нажмите кнопку 'Создать карточку'\n"+"t.me/poputka024",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: kbCreateTrip(slices.Contains(agentIds, user.ChatID), "user"),
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
	chatId := groupId

	lastMessageId, err := bot.mc.GetGroupLastMessageId()
	if err == nil {
		bot.tg.DeleteMessage(groupId, lastMessageId, &gotgbot.DeleteMessageOpts{})
	}

	// publish card to group
	cardMessage, err := bot.tg.SendPhoto(
		chatId,
		gotgbot.InputFileByReader("img.jpg", bytes.NewReader(card)),
		&gotgbot.SendPhotoOpts{
			ReplyMarkup: kbUnderCardInGroup(chat, trip),
		},
	)
	if err != nil {
		return nil, err
	}

	// update last group message
	groupLastMessage, err := bot.tg.SendMessage(
		chatId,
		"Что вы хотите сделать? ☺️",
		&gotgbot.SendMessageOpts{
			DisableNotification: true,
			ReplyMarkup:         kbOpenBot(),
		},
	)
	if err != nil {
		return nil, err
	}
	bot.mc.SetGroupLastMessageId(groupLastMessage.MessageId)

	return cardMessage, nil
}
