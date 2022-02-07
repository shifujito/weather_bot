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
		fmt.Print(err)
	}
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!!!!"
	fmt.Fprintf(w, msg)
	loadEnv()
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		fmt.Println("環境変数", err)
	}
	events, _ := bot.ParseRequest(r)
	fmt.Println("events", events)
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
	fmt.Println(r)
	fmt.Println(r.Body)
	events, _ := bot.ParseRequest(r)
	fmt.Println("events", events)
}

func main() {
	// http.HandleFunc("/", lineHandler)
	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: nil,
	}
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/callback", lineHandler)
	log.Fatal(server.ListenAndServe())
}
