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

type Config struct {
	db                  database.Config
	TgBotApiToken       string
	CreateTripWebAppUrl string
	SearchTripWebAppUrl string
}

func main() {
	log.SetOutput(os.Stdout)

	cfg := Config{
		db: database.Config{
			Host: os.Getenv("POSTGRES_HOST"),
			Port: os.Getenv("POSTGRES_PORT"),
			User: os.Getenv("POSTGRES_USER"),
			Pass: os.Getenv("POSTGRES_PASSWORD"),
			Name: os.Getenv("POSTGRES_DB"),
		},
		TgBotApiToken:       os.Getenv("TG_BOT_API_TOKEN"),
		CreateTripWebAppUrl: os.Getenv("CREATE_TRIP_WEB_APP_URL"),
		SearchTripWebAppUrl: os.Getenv("SEARCH_TRIP_WEB_APP_URL"),
	}

	database.Migrate(cfg.db)
	db, err := database.Initialize(cfg.db)
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
	bot, err := tg.NewBot(cfg.TgBotApiToken, &botOpts)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	ttlcache := ttlcache.New(ttlcache.WithTTL[int64, []model.Trip](time.Hour))

	longPoll := telegram.NewLongPoll(bot, db.Conn, ttlcache, cfg.CreateTripWebAppUrl, cfg.SearchTripWebAppUrl)
	go longPoll.Run()

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping the bot...")
}
