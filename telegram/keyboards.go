package telegram

import "github.com/PaulSonOfLars/gotgbot/v2"

func kbSelectRole() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Пассажир", CallbackData: "/user"},
			{Text: "Водитель", CallbackData: "/driver"},
		}},
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

func kbCreateTrip() gotgbot.ReplyKeyboardMarkup {
	return gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{{
			{Text: "Создать поездку", WebApp: &gotgbot.WebAppInfo{Url: createTrip}},
		}},
	}
}
