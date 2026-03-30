package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type GeocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type OpenMetioResponse struct {
	CurrentWeather struct {
		Temperature float32 `json:"temperature"`
		WindSpeed   float32 `json:"windspeed"`
		WeatherCode int     `json:"weathercode"`
	} `json:"current_weather"`
	CurrentWeatherUnits struct {
		Temperature string `json:"temperature"`
		WindSpeed   string `json:"windspeed"`
	} `json:"current_weather_units"`
}

type unitWeather struct {
	temperature string
	windSpeed   string
}

type weather struct {
	sity        string
	temperature float32
	windSpeed   float32
	weatherCode int
	unit        unitWeather
}

type cacheItem struct {
	weatherItem weather
	sity        string
	time        time.Time
}

type savedWeather struct {
	saveWeather map[string]cacheItem
	mu          sync.RWMutex
}

var cache savedWeather = savedWeather{
	saveWeather: map[string]cacheItem{},
	mu:          sync.RWMutex{},
}

func addToCash(weat weather, sity string) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.saveWeather[sity] = cacheItem{
		weatherItem: weat,
		sity:        sity,
		time:        time.Now(),
	}
	return nil
}

func readCache(sity string) error {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
}

func urlRequest(apiUrl string, params []any, target interface{}) (err error) {
	client := http.Client{Timeout: 30 * time.Second}
	url := fmt.Sprintf(apiUrl, params...)
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	if code := resp.StatusCode; code != 200 {
		err = errors.New("Код не 200.")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, target)
	if err != nil {
		return
	}
	return
}

func getWeather(sity string) (weather, error) {
	var geoResp GeocodingResponse
	err := urlRequest("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json", []any{sity}, &geoResp)
	if err != nil {
		return weather{}, err
	}
	if len(geoResp.Results) != 1 {
		return weather{}, errors.New("Город или не найден или не корректный.")
	}
	var metioResp OpenMetioResponse
	err = urlRequest(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true&timezone=auto",
		[]any{geoResp.Results[0].Latitude, geoResp.Results[0].Longitude},
		&metioResp,
	)
	if err != nil {

	}
	return weather{
		sity:        geoResp.Results[0].Name,
		temperature: metioResp.CurrentWeather.Temperature,
		windSpeed:   metioResp.CurrentWeather.WindSpeed,
		weatherCode: metioResp.CurrentWeather.WeatherCode,
		unit: unitWeather{
			temperature: metioResp.CurrentWeatherUnits.Temperature,
			windSpeed:   metioResp.CurrentWeatherUnits.WindSpeed,
		},
	}, nil
}

func outWeather(weat weather) {
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Город: %s\n", weat.sity)
	fmt.Printf("Температура: %v%s\n", weat.temperature, weat.unit.temperature)
	fmt.Printf("Скорость ветра: %v%s\n", weat.windSpeed, weat.unit.windSpeed)
	fmt.Printf("Состояние: %s%s\n", codeWeather(weat.weatherCode), getWeatherIcon(weat.weatherCode))
	fmt.Println("-------------------------------------------------")
}

func codeWeather(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Облачно"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 80, 81, 82:
		return "Ливень"
	case 95, 96, 99:
		return "Гроза"
	default:
		return "Неизвестно"
	}
}

func getWeatherIcon(code int) string {
	switch code {
	case 0:
		return "☀️"
	case 1:
		return "🌤️"
	case 2:
		return "⛅"
	case 3:
		return "☁️"
	case 45, 48:
		return "🌫️"
	case 51, 53, 55:
		return "🌧️"
	case 56, 57:
		return "🌧️❄️"
	case 61, 63, 65:
		return "🌧️"
	case 66, 67:
		return "🌧️❄️"
	case 71, 73, 75:
		return "🌨️"
	case 77:
		return "❄️"
	case 80, 81, 82:
		return "🌧️☔"
	case 85, 86:
		return "🌨️☔"
	case 95, 96, 99:
		return "⛈️"
	default:
		return "❓"
	}
}

func funcGetWeather() {
	for {
		var sity string
		_, err := fmt.Fscan(os.Stdin, &sity)
		if err != nil {
			fmt.Println("[ERROR]: ", err)
			continue
		}
		if sity == "exit" {
			return
		}
		weat, err := getWeather(sity)
		if err != nil {
			fmt.Println("[ERROR]: ", err)
			continue
		}

		// fmt.Println(weat)
		outWeather(weat)
	}

}

func main() {
	funcGetWeather()
}
