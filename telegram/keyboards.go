package telegram

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/p24/model"
)

func kbSelectRole() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Пассажир", CallbackData: "/user"},
			{Text: "Водитель", CallbackData: "/driver"},
		}},
	}
}

func kbUnderCardInGroup(chat *gotgbot.Chat, trip *model.Trip) gotgbot.InlineKeyboardMarkup {
	button := gotgbot.InlineKeyboardButton{Text: trip.Phone, CopyText: &gotgbot.CopyTextButton{Text: trip.Phone}}
	if trip.Phone == "" {
		button = gotgbot.InlineKeyboardButton{Text: "Написать в ЛС", Url: "t.me/" + chat.Username}
	}

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{button}},
	}
}

func kbOpenBot() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{Text: "Написать боту", Url: "t.me/poputka24bot?start=hello"}},
			{{Text: "Поиск", Url: "t.me/poputka24bot/search"}},
		},
	}
}

func kbOpenGroup() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Перейти в группу", Url: "t.me/poputka024"},
		}},
	}
}

func kbUnderCard(trip *model.Trip) gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{Text: "Перейти в группу️", Url: fmt.Sprintf("t.me/poputka024/%d", trip.MessageId)}},
		},
	}
}

func kbSendContact() gotgbot.ReplyKeyboardMarkup {
	return gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{{
			{Text: "Отправить контакт", RequestContact: true},
		}},
	}
}

func kbFindPassengers() gotgbot.ReplyKeyboardMarkup {
	return gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{{
			{Text: "Найти попутчиков", WebApp: &gotgbot.WebAppInfo{Url: searchTrips}},
		}},
	}
}

func kbCreateTrip(admin bool, userType string) gotgbot.ReplyKeyboardMarkup {
	link := createTrip + "?user_type=" + userType
	if admin {
		link = createTripAdmin
	}
	return gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{{
			{Text: "Создать карточку", WebApp: &gotgbot.WebAppInfo{Url: link}},
		}},
	}
}
