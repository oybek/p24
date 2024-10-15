package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/greetings/telegram"
	tg "github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	log.SetOutput(os.Stdout)

	//
	tgBotApiToken := os.Getenv("TG_BOT_API_TOKEN")

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
