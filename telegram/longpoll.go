package telegram

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/google/uuid"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/model"
)

type LongPoll struct {
	bot *gotgbot.Bot
	db  *sql.DB
}

func NewLongPoll(
	bot *gotgbot.Bot,
	db *sql.DB,
) *LongPoll {
	return &LongPoll{
		bot: bot,
		db:  db,
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
	dispatcher.AddHandler(handlers.NewMessage(
		func(msg *gotgbot.Message) bool { return strings.HasPrefix(msg.Text, "/start") },
		lp.handleStart,
	))

	dispatcher.AddHandler(handlers.NewMessage(message.Text, lp.handleText))

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

func (lp *LongPoll) handleStart(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	text := ctx.EffectiveMessage.Text

	UUID, err := uuid.Parse(strings.TrimPrefix(text, "/start "))
	if err != nil {
		log.Printf("error parsing uuid from '%s': %s", text, err.Error())
		return err
	}

	user := model.User{
		ChatId: chat.Id,
		UUID:   UUID,
	}

	_, err = database.Transact(lp.db, func(tx database.TransactionOps) (any, error) {
		return database.UpsertUser(tx, &user)
	})
	if err != nil {
		log.Printf("error upserting user %#v: %s", user, err.Error())
		return err
	}

	return nil
}

func (lp *LongPoll) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	text := ctx.EffectiveMessage.Text

	_, err := b.SendMessage(chat.Id, text, &gotgbot.SendMessageOpts{})
	return err
}
