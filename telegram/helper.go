package telegram

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/choguuket/model"
)

func parse[T model.Validated](jsonRaw string) (*T, error) {
	var data T
	if err := json.Unmarshal([]byte(jsonRaw), &data); err != nil {
		return nil, err
	}
	if !data.IsValid() {
		return nil, errors.New("invalid data")
	}
	return &data, nil
}

func (lp *LongPoll) sendText(chatId int64, text string) error {
	_, err := lp.bot.SendMessage(chatId, text, &gotgbot.SendMessageOpts{
		ParseMode: "markdown",
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{Text: createTripButtonText, WebApp: &gotgbot.WebAppInfo{Url: fmt.Sprintf("%s?chatId=%d", lp.createTripWebAppUrl, chatId)}},
					{Text: searchTripButtonText, WebApp: &gotgbot.WebAppInfo{Url: fmt.Sprintf("%s?chatId=%d", lp.searchTripWebAppUrl, chatId)}}}},
			ResizeKeyboard: true,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (lp *LongPoll) sendTrip(chat *gotgbot.Chat, tripReqId int64, trips []model.Trip) error {
	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "markdown",
	}

	if len(trips) > 1 {
		sendMessageOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				gotgbot.InlineKeyboardButton{
					Text:         fmt.Sprintf("Еще %d", len(trips)-1),
					CallbackData: fmt.Sprintf("next;%d", tripReqId),
				},
			}},
		}
	}

	_, err := lp.bot.SendMessage(chat.Id, trips[0].String(), &sendMessageOpts)
	return err
}

func (lp *LongPoll) sendTripNotification(
	tripReq *model.TripReq,
	trip *model.Trip,
) error {
	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "markdown",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				gotgbot.InlineKeyboardButton{
					Text:         "Больше не нужно",
					CallbackData: fmt.Sprintf("del;%d", 1),
				},
			}},
		},
	}

	_, err := lp.bot.SendMessage(
		tripReq.ChatId,
		fmt.Sprintf(
			"*По Вашему запросу:*\n%s\n\n*Новая поездка:*\n%s",
			tripReq.String(),
			trip.String(),
		),
		&sendMessageOpts,
	)

	return err
}
