package services

import (
	"fmt"
	"net/http"
)

type WeatherService interface {
	GetWeather(city, country string) (*http.Response, error)
}

type OpenWeatherService struct {
	URL string
}

type WeatherFromFileService struct {
	Route string
}

func (o OpenWeatherService) GetWeather(city, country string) (*http.Response, error) {

	cityCode := city + "," + country
	url := fmt.Sprintf(o.URL, cityCode)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (wf WeatherFromFileService) GetWeather(city, country string) (*http.Response, error) {
	return nil, nil
}
