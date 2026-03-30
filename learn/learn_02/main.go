

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"
// )

// type GeocodingResponse struct {
// 	Results []struct {
// 		Name      string  `json:"name"`
// 		Latitude  float64 `json:"latitude"`
// 		Longitude float64 `json:"longitude"`
// 	} `json:"results"`
// }

// type OpenMetioResponse struct {
// 	CurrentWeather struct {
// 		Temperature float32 `json:"temperature"`
// 		WindSpeed   float32 `json:"windspeed"`
// 		WeatherCode int     `json:"weathercode"`
// 	} `json:"current_weather"`
// }

// type weather struct {
// 	sity        string
// 	temperature float32
// 	windSpeed   float32
// 	weatherCode int
// }

// func urlRequest(apiUrl string, params []any, target interface{}) (err error) {
// 	client := http.Client{Timeout: 30 * time.Second}
// 	url := fmt.Sprintf(apiUrl, params...)
// 	resp, err := client.Get(url)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	if code := resp.StatusCode; code != 200 {
// 		fmt.Println(code)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	err = json.Unmarshal(body, target)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	return
// }

// func getWeather(sity string) weather {
// 	var geoResp GeocodingResponse
// 	err := urlRequest("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json", []any{sity}, &geoResp)
// 	if err != nil {
// 		// return
// 	}
// 	var metioResp OpenMetioResponse
// 	err = urlRequest(
// 		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true&timezone=auto",
// 		[]any{geoResp.Results[0].Latitude, geoResp.Results[0].Longitude},
// 		&metioResp,
// 	)
// 	return weather{
// 		sity:        geoResp.Results[0].Name,
// 		temperature: metioResp.CurrentWeather.Temperature,
// 		windSpeed:   metioResp.CurrentWeather.WindSpeed,
// 		weatherCode: metioResp.CurrentWeather.WeatherCode,
// 	}
// }

// func main() {
// 	weat := getWeather("Moscow")
// 	fmt.Println(weat)
// }
