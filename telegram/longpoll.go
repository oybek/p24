package telegram

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/sashabaranov/go-openai"
)

type LongPoll struct {
	bot          *gotgbot.Bot
	db           *sql.DB
	openaiClient *openai.Client
}

func NewLongPoll(
	bot *gotgbot.Bot,
	db *sql.DB,
	openaiClient *openai.Client,
) *LongPoll {
	return &LongPoll{
		bot:          bot,
		db:           db,
		openaiClient: openaiClient,
	}
}

const createAptekaWebAppUrl = "https://wolfrepos.github.io/apteka/create/index.html"

func (lp *LongPoll) Run() {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewMessage(
		func(msg *gotgbot.Message) bool { return strings.HasPrefix(msg.Text, "/create_apteka") },
		lp.handleCreateApteka,
	))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, lp.handleText))
	dispatcher.AddHandler(handlers.NewMessage(message.Voice, lp.handleVoice))

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

func (lp *LongPoll) handleVoice(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	voice := ctx.EffectiveMessage.Voice

	if voice.Duration > 20 {
		return lp.sendText(chat.Id, TextTooLongVoice)
	}

	file, err := b.GetFile(voice.FileId, &gotgbot.GetFileOpts{})
	if err != nil {
		return err
	}

	resp, err := http.Get(file.URL(b, &gotgbot.RequestOpts{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		Reader:   resp.Body,
		FilePath: file.FilePath,
		Prompt:   "Парацетамол, ТайлолХот, Тримол",
		Language: "ru",
	}

	context := context.Background()
	openaiResp, err := lp.openaiClient.CreateTranscription(context, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return err
	}

	return lp.sendText(chat.Id, openaiResp.Text)
}

func (lp *LongPoll) handleText(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	return lp.sendText(chat.Id, TextDefault)
}

func (lp *LongPoll) sendText(chatId int64, text string) error {
	_, err := lp.bot.SendMessage(chatId, text, &gotgbot.SendMessageOpts{})
	return err
}

func (lp *LongPoll) handleCreateApteka(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveMessage.Chat
	createAptekaKeyboard := &gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: "Создать аптеку", WebApp: &gotgbot.WebAppInfo{Url: createAptekaWebAppUrl}},
			},
		},
	}
	_, err := lp.bot.SendMessage(chat.Id, TextCreateApteka,
		&gotgbot.SendMessageOpts{ReplyMarkup: createAptekaKeyboard})
	return err
}
