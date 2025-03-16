package telegram

import (
	"embed"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/oybek/p24/mongo"
	"github.com/oybek/p24/tools"
)

type Bot struct {
	tg        *gotgbot.Bot
	mc        *mongo.MongoClient
	fonts     *embed.FS
	cityNames *tools.BMap[string, string]
}

func NewBot(
	tg *gotgbot.Bot,
	mc *mongo.MongoClient,
	fonts *embed.FS,
	cityNames *tools.BMap[string, string],
) *Bot {
	return &Bot{
		tg:        tg,
		mc:        mc,
		fonts:     fonts,
		cityNames: cityNames,
	}
}

const searchTrips = "https://oybek.github.io/p24-wa/?user_type=driver"
const createTrip = "https://oybek.github.io/p24-wa/?user_type=user"
const createTripAdmin = "https://oybek.github.io/p24-wa/?user_type=admin"

var agentIds = []int64{108683062}

func (lp *Bot) Run() {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// admin commands
	dispatcher.AddHandler(handlers.NewMessage(message.HasPrefix("/new_city"), lp.handleCommandNewCity))

	//
	dispatcher.AddHandler(handlers.NewMessage(message.HasPrefix("/group_update"), lp.handleCommandGroupUpdate))
	dispatcher.AddHandler(handlers.NewMessage(message.HasPrefix("/profile"), lp.handleCommandProfile))
	dispatcher.AddHandler(handlers.NewMessage(message.HasPrefix("/change"), lp.handleCommandChange))
	dispatcher.AddHandler(handlers.NewMessage(message.HasPrefix("/start"), lp.handleStartCommand))
	dispatcher.AddHandler(handlers.NewMessage(message.HasPrefix("/help"), lp.handleCommandHelp))
	dispatcher.AddHandler(handlers.NewMessage(message.Contact, lp.handleContact))
	dispatcher.AddHandler(handlers.NewMessage(message.Photo, lp.handlePhoto))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, lp.handleText))
	dispatcher.AddHandler(handlers.NewMessage(messageWebApp, lp.handleWebAppData))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("/driver"), lp.handleCommandDriver))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("/user"), lp.handleCommandUser))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("/show_phone"), lp.handleCommandShowPhone))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("/del"), lp.handleCommandDeleteTrip))

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
