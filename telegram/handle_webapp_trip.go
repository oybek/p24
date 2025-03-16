package telegram

import (
	"bytes"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/p24/model"
)

func (bot *Bot) handleWebAppTrip(chat *gotgbot.Chat, trip *model.Trip) error {
	trip.ChatID = chat.Id
	trip.State = "active"

	tripID, err := bot.mc.TripCreate(trip)
	if err != nil {
		return err
	}
	trip.ID = tripID

	user, err := bot.mc.UserGetByChatID(trip.ChatID)
	if err != nil {
		return err
	}

	tripView := MapToTripView(trip, user)
	font, _ := bot.fonts.ReadFile("fonts/lcd5x8h.ttf")
	tripCard, err := DrawTextToImage(FormatTrip(tripView), font)
	if err != nil {
		return err
	}

	messageInGroup, err := bot.tg.SendPhoto(
		groupId,
		gotgbot.InputFileByReader("img.jpg", bytes.NewReader(tripCard)),
		&gotgbot.SendPhotoOpts{
			ReplyMarkup: kbUnderCardInGroup(chat, trip),
		},
	)
	if err != nil {
		return err
	}

	trip.MessageId = messageInGroup.MessageId
	log.Printf("TripUpdateMessageID: %s %d\n", trip.ID, trip.MessageId)
	err = bot.mc.TripUpdateMessageID(trip.ID, trip.MessageId)
	if err != nil {
		return err
	}

	_, err = bot.tg.SendPhoto(
		chat.Id,
		gotgbot.InputFileByReader("img.jpg", bytes.NewReader(tripCard)),
		&gotgbot.SendPhotoOpts{
			Caption: "✅ Ваша карточка успешно создана! В ближайшее время с вами свяжутся наши водители.\n\n" +
				"❗ Как только вы договоритесь с водителем, пожалуйста, удалите карточку, нажав кнопку ниже.",
			ReplyMarkup: kbUnderCard(trip),
		},
	)
	if err != nil {
		return err
	}

	return err
}
