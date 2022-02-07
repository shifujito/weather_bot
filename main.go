package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func loadEnv() (bot *linebot.Client) {
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := linebot.NewTextMessage("hello, world")
	bot := loadEnv()
	// テキストメッセージを友達登録しているユーザー全員に配信する
	if _, err := bot.BroadcastMessage(msg).Do(); err != nil {
		log.Fatal(err)
	}
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	bot := loadEnv()
	events, _ := bot.ParseRequest(r)
	for _, event := range events {
		fmt.Println(event.Message)
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func main() {
	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: nil,
	}
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/callback", lineHandler)
	log.Fatal(server.ListenAndServe())
}
