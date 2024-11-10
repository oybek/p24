package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	users, err := database.Transact(lp.db, func(tx database.TransactionOps) ([]model.User, error) {
		return database.SelectUser(tx, UUID)
	})
	if err != nil {
		log.Printf("error selecting user: %s", err.Error())
		return err
	}

	if len(users) > 0 {
		_, err = b.SendMessage(chat.Id,
			"Регистрация неуспешна! 😔\n"+
				"По данному QR уже была регистрация",
			&gotgbot.SendMessageOpts{})
		if err != nil {
			log.Printf("error sending a message: %s", err.Error())
			return err
		}
		return nil
	}

	user := model.User{
		ChatId: chat.Id,
		UUID:   UUID,
	}
	_, err = database.Transact(lp.db, func(tx database.TransactionOps) (any, error) {
		return database.UpsertUser(tx, &user)
	})
	if err != nil {
		log.Printf("error inserting user %#v: %s", user, err.Error())
		return err
	}

	_, err = b.SendMessage(chat.Id,
		"Регистрация успешна! ✅\n"+
			"Теперь наклейте второй QR код на Вашу машину так чтобы другие могли ее отсканировать 😊",
		&gotgbot.SendMessageOpts{})
	if err != nil {
		log.Printf("error sending a message: %s", err.Error())
		return nil
	}

	return nil
}

func (lp *LongPoll) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	text := ctx.EffectiveMessage.Text

	_, err := b.SendMessage(chat.Id, text, &gotgbot.SendMessageOpts{})
	return err
}

func (lp *LongPoll) NotifyUser(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	UUID, err := uuid.Parse(query.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := database.Transact(lp.db, func(tx database.TransactionOps) ([]model.User, error) {
		return database.SelectUser(tx, UUID)
	})
	if err != nil {
		log.Printf("error selecting user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = lp.bot.SendMessage(users[0].ChatId, "Вас просят переставить машину!", &gotgbot.SendMessageOpts{})
	if err != nil {
		log.Printf("error sending a message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
