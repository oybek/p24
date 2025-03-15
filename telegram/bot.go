package telegram

import (
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/oybek/p24/mongo"
)

type Bot struct {
	tg *gotgbot.Bot
	mc *mongo.MongoClient
}

func NewBot(
	tg *gotgbot.Bot,
	mc *mongo.MongoClient,
) *Bot {
	return &Bot{
		tg: tg,
		mc: mc,
	}
}

const searchTrips = "https://oybek.github.io/p24-wa/?user_type=driver"
const createTrip = "https://oybek.github.io/p24-wa/?user_type=user"

func (lp *Bot) Run() {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewMessage(message.Command, lp.handleCommand))
	dispatcher.AddHandler(handlers.NewMessage(message.Contact, lp.handleContact))
	dispatcher.AddHandler(handlers.NewMessage(message.Photo, lp.handlePhoto))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, lp.handleText))
	dispatcher.AddHandler(handlers.NewMessage(messageWebApp, lp.handleWebAppData))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("/driver"), lp.handleCommandDriver))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("/user"), lp.handleCommandUser))

	// Start receiving updates.
	err := updater.StartPolling(lp.tg, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	lp.SetupCommands()

	log.Printf("%s has been started...\n", lp.tg.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func messageWebApp(msg *gotgbot.Message) bool {
	return msg.WebAppData != nil
}
