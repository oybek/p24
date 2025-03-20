package telegram

import (
	"bytes"
	"image/color"
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

	tripView := bot.MapToTripView(trip, user)
	font, _ := bot.fonts.ReadFile("fonts/lcd5x8h.ttf")

	cardColor := color.RGBA{R: 200, G: 250, B: 200, A: 255}
	if user.UserType == "user" {
		cardColor = color.RGBA{R: 250, G: 200, B: 200, A: 255}
	}
	tripCard, err := DrawTextToImage(FormatTrip(tripView, user.UserType), font, cardColor)
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
			Caption:     "✅ Ваша карточка успешно добавлена в группу!",
			ReplyMarkup: kbUnderCard(trip),
		},
	)
	if err != nil {
		return err
	}

	return err
}
