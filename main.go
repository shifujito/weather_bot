package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

func GetWeather(code string) string {
	json := httpGetStr("https://www.jma.go.jp/bosai/forecast/data/overview_forecast/" + code + "0000.json")
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
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				fmt.Println(message.Text, event.ReplyToken)
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

func getPrefectureCode(text string) (string, bool) {
	pref_code_map := map[string]string{"北海道": "02", "青森": "02", "岩手": "03", "宮城": "04", "秋田": "05", "山形": "06", "福島": "07",
		"茨城": "08", "栃木": "09", "群馬": "10", "埼玉": "11", "千葉": "12", "東京": "13", "神奈川": "14", "新潟": "15", "富山": "16", "石川": "17",
		"福井": "18", "山梨": "19", "長野": "20", "岐阜": "21", "静岡": "22", "愛知": "23", "三重": "24", "滋賀": "25", "京都": "26", "大阪": "27",
		"兵庫": "28", "奈良": "29", "和歌山": "30", "鳥取": "31", "島根": "32", "岡山": "33", "広島": "34", "山口": "35", "徳島": "36", "香川": "37", "愛媛": "38", "高知": "39",
		"福岡": "40", "佐賀": "41", "長崎": "42", "熊本": "43", "大分": "44", "宮崎": "45", "鹿児島": "45", "沖縄": "45"}
	replaced_text := strings.Replace(text, "都", "", -1)
	replaced_text = strings.Replace(replaced_text, "県", "", -1)
	replaced_text = strings.Replace(replaced_text, "府", "", -1)
	code, hasCode := pref_code_map[replaced_text]
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
