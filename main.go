package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tg "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/jub0bs/fcors"
	"github.com/oybek/p24/mongo"
	"github.com/oybek/p24/rest"
	"github.com/oybek/p24/telegram"
	"github.com/oybek/p24/tools"
)

type Config struct {
	mongoURL    string
	botAPIToken string
}

//go:embed fonts/*
var fonts embed.FS

func main() {
	ctx := context.Background()

	//
	log.SetOutput(os.Stdout)

	//
	cfg := Config{
		mongoURL:    os.Getenv("MONGO_URL"),
		botAPIToken: os.Getenv("BOT_API_TOKEN"),
	}

	//
	tgbot, err := tg.NewBot(cfg.botAPIToken, &tg.BotOpts{
		BotClient: &tg.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &tg.RequestOpts{
				Timeout: 10 * time.Second,
				APIURL:  tg.DefaultAPIURL,
			},
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	//
	mc, err := mongo.NewMongoClient(ctx, cfg.mongoURL)
	if err != nil {
		panic("failed to create new mongo client: " + err.Error())
	}

	//
	bmap := tools.NewBMap[string, string]()

	bot := telegram.NewBot(tgbot, mc, &fonts, bmap)
	err = bot.InitCityNames()
	if err != nil {
		panic("can't load city names from mongo: " + err.Error())
	}

	go bot.Run()

	cors, _ := fcors.AllowAccess(
		fcors.FromAnyOrigin(),
		fcors.WithMethods(
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		),
		fcors.WithRequestHeaders("Authorization"),
	)

	r := rest.New(bot, mc)
	http.Handle("/ok", cors(http.HandlerFunc(r.Ok)))
	http.Handle("/trips", cors(http.HandlerFunc(r.TripFind)))
	http.Handle("/cards", cors(http.HandlerFunc(r.TripCard)))
	http.Handle("/cities", cors(http.HandlerFunc(r.Cities)))
	go http.ListenAndServe(":5555", nil)

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping the bot...")
}
