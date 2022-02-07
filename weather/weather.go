package weather

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type WeatherType struct {
	Area string `json:"targetArea"`
	Body string `json:"text"`
}

func GetWeather() string {
	json := httpGetStr("https://www.jma.go.jp/bosai/forecast/data/overview_forecast/130000.json")
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
