package telegram

import (
	"bytes"
	"image/color"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/p24/model"
)

func (bot *Bot) handleWebAppTrip(chat *gotgbot.Chat, trip *model.Trip) error {
	user, err := bot.mc.UserGetByChatID(chat.Id)
	if err != nil {
		return err
	}

	trip.ChatID = chat.Id
	trip.UserType = user.UserType

	tripID, err := bot.mc.TripCreate(trip)
	if err != nil {
		return err
	}
	trip.ID = tripID

	tripCard, err := bot.DrawCard(trip, user)
	if err != nil {
		return err
	}

	messageInGroup, err := bot.publishCard(chat, trip, tripCard)
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
			Caption: "✅ Ваша карточка успешно добавлена в группу!\n\n" +
				"☝️ Добавьте в группу ещё трёх человек, чтобы создавать больше карточек",
			ReplyMarkup: kbUnderCard(trip),
		},
	)
	if err != nil {
		return err
	}

	return err
}

func (bot *Bot) DrawCard(trip *model.Trip, user *model.User) ([]byte, error) {
	tripView := bot.MapToTripView(trip, user)
	font, _ := bot.fonts.ReadFile("fonts/lcd5x8h.ttf")

	cardColor := color.RGBA{R: 200, G: 250, B: 200, A: 255}
	if trip.UserType == "passenger" {
		cardColor = color.RGBA{R: 250, G: 200, B: 200, A: 255}
	}
	tripCard, err := DrawTextToImage(FormatTrip(tripView, user.UserType), font, cardColor)
	if err != nil {
		return nil, err
	}
	return tripCard, nil
}
