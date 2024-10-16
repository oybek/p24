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
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/database"
	"github.com/oybek/choguuket/model"
	"github.com/oybek/choguuket/telegram"
)

func main() {
	log.SetOutput(os.Stdout)

	dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		tgBotApiToken,
		createTripWebAppUrl,
		searchTripWebAppUrl :=
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("TG_BOT_API_TOKEN"),
		os.Getenv("CREATE_TRIP_WEB_APP_URL"),
		os.Getenv("SEARCH_TRIP_WEB_APP_URL")

	database.Migrate(dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := database.Initialize(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	defer db.Conn.Close()

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

	ttlcache := ttlcache.New(ttlcache.WithTTL[int64, []model.Trip](time.Hour))

	longPoll := telegram.NewLongPoll(bot, db.Conn, ttlcache, createTripWebAppUrl, searchTripWebAppUrl)
	go longPoll.Run()

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping the bot...")
}
