package weather

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetWeather() string {
	json := httpGetStr("https://www.jma.go.jp/bosai/forecast/data/overview_forecast/130000.json")
	return json
}

func httpGetStr(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Get Http Error", err)
	}
	body, _ := io.ReadAll(response.Body)
	return string(body)

}

func Weather() {
	result := GetWeather()
	fmt.Println(result)
}
