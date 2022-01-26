package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("環境変数が読み込めませんでした。")
	}
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	loadEnv()
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	message := linebot.NewTextMessage("hello, world")
	// テキストメッセージを友達登録しているユーザー全員に配信する
	if _, err := bot.BroadcastMessage(message).Do(); err != nil {
		log.Fatal(err)

	}

}

func main() {
	// http.HandleFunc("/", helloHandler)
	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: nil,
	}
	http.HandleFunc("/", lineHandler)
	log.Fatal(server.ListenAndServe())
}
