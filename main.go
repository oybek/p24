package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tg "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/telegram"
)

func main() {
	log.SetOutput(os.Stdout)

	dbHost, dbPort, dbUser, dbPassword, dbName, tgBotApiToken, _, _ :=
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("TG_BOT_API_TOKEN"),
		os.Getenv("CREATE_TRIP_WEB_APP_URL"),
		os.Getenv("SEARCH_TRIP_WEB_APP_URL")

	database.Migrate(dbHost, dbPort, dbUser, dbPassword, dbName)

	//
	botOpts := tg.BotOpts{
		BotClient: &tg.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &tg.RequestOpts{
				Timeout: 10 * time.Second,
				APIURL:  tg.DefaultAPIURL,
			},
		},
	}
	bot, err := tg.NewBot(tgBotApiToken, &botOpts)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	longPoll := telegram.NewLongPoll(bot)
	go longPoll.Run()

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping the bot...")
}
