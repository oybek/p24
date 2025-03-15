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

func kbShowPhone() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Показать номер", CallbackData: "/show_phone123123"},
		}},
	}
}

func bkOpenBot() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Открыть", Url: "t.me/poputka24bot"},
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
