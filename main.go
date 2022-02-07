package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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

type WeatherType struct {
	Area string `json:"targetArea"`
	Body string `json:"text"`
}

func GetWeather(code int) string {
	strCode := strconv.Itoa(code)
	json := httpGetStr("https://www.jma.go.jp/bosai/forecast/data/overview_forecast/" + strCode + "0000.json")
	weather := formatWeather(json)
	return weather.Area + weather.Body
}

func httpGetStr(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Get Http Error", err)
	}
	body, _ := io.ReadAll(response.Body)
	return string(body)

}

func formatWeather(str string) *WeatherType {
	weather := new(WeatherType)
	if err := json.Unmarshal([]byte(str), weather); err != nil {
		log.Fatal("JSON Unmarshal error:", err)
	}
	return weather
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
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(textParser(message.Text))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func textParser(text string) string {
	code, hasCode := getPrefectureCode(text)
	if !hasCode {
		return "天気予報が知りたい都道府県を入力してください。"
	}
	result := GetWeather(code)
	return result
}

func getPrefectureCode(text string) (int, bool) {
	pref_code_map := map[string]int{"東京都": 13}
	code, hasCode := pref_code_map[text]
	return code, hasCode
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
