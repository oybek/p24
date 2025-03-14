package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tg "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gorilla/mux"
	"github.com/jub0bs/fcors"
	"github.com/oybek/p24/mongo"
	"github.com/oybek/p24/telegram"
)

type Config struct {
	mongoURL    string
	botAPIToken string
}

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

	bot := telegram.NewBot(tgbot, mc)
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

	r := mux.NewRouter()
	http.Handle("/", cors(r))
	go http.ListenAndServe(":5556", nil)

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping the bot...")
}
