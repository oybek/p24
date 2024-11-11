package telegram

import (
	"database/sql"
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
		return err
	}

	users, err := database.Transact(lp.db, func(tx database.TransactionOps) ([]model.User, error) {
		return database.SelectUser(tx, UUID)
	})
	if err != nil {
		return err
	}

	if len(users) > 0 {
		if err = lp.sendText(chat.Id, TextWhenFailStart); err != nil {
			return err
		}
		return nil
	}

	chatInfo, err := lp.bot.GetChat(chat.Id, &gotgbot.GetChatOpts{})
	if err != nil {
		return err
	}

	user := model.User{
		ChatId: chat.Id,
		UUID:   UUID,
		Nick:   chatInfo.Username,
	}
	_, err = database.Transact(lp.db, func(tx database.TransactionOps) (any, error) {
		return database.UpsertUser(tx, &user)
	})
	if err != nil {
		return err
	}

	if err = lp.sendText(chat.Id, TextWhenOkStart); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	if err = lp.sendText(chat.Id, TextMoveCar); err != nil {
		return nil
	}

	return nil
}

func (lp *LongPoll) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	return lp.sendText(chat.Id, TextDefault)
}

func (lp *LongPoll) NotifyUser(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) < 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = lp.sendText(users[0].ChatId, TextMoveCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (lp *LongPoll) sendText(chatId int64, text string) error {
	_, err := lp.bot.SendMessage(chatId, text, &gotgbot.SendMessageOpts{})
	return err
}
