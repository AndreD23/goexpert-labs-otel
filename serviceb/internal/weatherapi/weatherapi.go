package weatherapi

import (
	"fmt"
	"github.com/AndreD23/goexpert-labs-otel/serviceb/pkg/utils"
	"net/url"
)

type WeatherAPI struct {
	APIKey string
}

type Response struct {
	Temperature struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		TempK float64 `json:"temp_k"`
	} `json:"current"`
}

func NewWeatherAPI(apiKey string) *WeatherAPI {
	return &WeatherAPI{
		APIKey: apiKey,
	}
}

func (w *WeatherAPI) GetTempByCity(city string) (Response, error) {
	wUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", w.APIKey, url.QueryEscape(city))
	var data Response

	err := utils.FetchData(wUrl, &data)
	if err != nil {
		return data, err
	}

	data.Temperature.TempK = data.Temperature.TempC + 273

	return data, nil
}
