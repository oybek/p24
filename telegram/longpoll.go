package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/model"
)

type LongPoll struct {
	bot         *gotgbot.Bot
	db          *sql.DB
	searchCache *ttlcache.Cache[int64, []model.Trip]
}

func NewLongPoll(
	bot *gotgbot.Bot,
	db *sql.DB,
	searchCache *ttlcache.Cache[int64, []model.Trip],
) *LongPoll {
	return &LongPoll{
		bot:         bot,
		db:          db,
		searchCache: searchCache,
	}
}

func (lp *LongPoll) Run() {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	//
	dispatcher.AddHandler(handlers.NewMessage(message.Text, lp.handleText))
	dispatcher.AddHandler(handlers.NewMessage(isWebAppData, lp.handleWebAppData))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("next"), lp.handleNextTrip))

	// Start receiving updates.
	err := updater.StartPolling(lp.bot, &ext.PollingOpts{
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

	log.Printf("%s has been started...\n", lp.bot.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func (lp *LongPoll) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	text := ctx.EffectiveMessage.Text

	if strings.HasPrefix(text, "/webapp") {
		context := ext.Context{}
		rawJson := strings.TrimSpace(strings.TrimPrefix(text, "/webapp"))
		context.EffectiveMessage = &gotgbot.Message{
			Chat:       chat,
			WebAppData: &gotgbot.WebAppData{Data: rawJson},
		}
		return lp.handleWebAppData(b, &context)
	}

	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

func isWebAppData(msg *gotgbot.Message) bool {
	return msg.WebAppData != nil
}
