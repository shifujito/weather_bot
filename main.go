package main

import (
	"fmt"
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

func main() {
	loadEnv()
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
}
