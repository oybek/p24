package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/p24/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func kbSelectRole() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Пассажир", CallbackData: "/user"},
			{Text: "Водитель", CallbackData: "/driver"},
		}},
	}
}

func kbShowPhone(tripId primitive.ObjectID) gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Показать номер", CallbackData: "/show_phone" + tripId.Hex()},
		}},
	}
}

func kbOpenBot() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Открыть", Url: "t.me/poputka24bot?start=hello"},
		}},
	}
}

func kbOpenGroup() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Перейти в группу", Url: "t.me/poputka24ads"},
		}},
	}
}

func kbUnderCard(trip *model.Trip) gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			//{{Text: "Перейти в группу ➡️", Url: fmt.Sprintf("t.me/poputka24ads/%d", trip.MessageId)}},
			{{Text: "Удалить карточку 🔥", CallbackData: "/del" + trip.ID.Hex()}},
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

func kbCreateTrip() gotgbot.ReplyKeyboardMarkup {
	return gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{{
			{Text: "Создать поездку", WebApp: &gotgbot.WebAppInfo{Url: createTrip}},
		}},
	}
}
