package models

import (
	"fmt"
	"strconv"
	"time"
)

type WeatherJsonData struct {
	Weather     []map[string]string    `json:"weather"`
	Main        map[string]float64     `json:"main"`
	Wind        map[string]float64     `json:"wind"`
	Sys         map[string]interface{} `json:"sys"`
	City        string                 `json:"name"`
	Coordinates map[string]float64     `json:"coord"`
}

type WeatherData struct {
	Location      string
	Temperature   string
	Wind          map[string]float64
	Cloudines     string
	Presure       string
	Humidity      string
	Sunrise       float64
	Sunset        float64
	Coordinates   map[string]float64
	RequestedTime time.Time `json:"Requested_time"`
}

func (w *WeatherData) GetDataFromJSON(data WeatherJsonData) error {

	country, ok := data.Sys["country"].(string)
	if !ok {
		err := fmt.Errorf("Country is not type string")
		return err
	}
	w.Location = data.City + ", " + country
	w.Temperature = strconv.FormatFloat(data.Main["temp"], 'f', 0, 64) + "°C"
	w.Wind = data.Wind
	w.Cloudines = data.Weather[0]["description"]
	w.Presure = strconv.FormatFloat(data.Main["pressure"], 'f', 0, 64) + " hpa"
	w.Humidity = strconv.FormatFloat(data.Main["humidity"], 'f', 0, 64) + "%"
	w.Sunrise = data.Sys["sunrise"].(float64)
	w.Sunset = data.Sys["sunset"].(float64)
	w.Coordinates = data.Coordinates
	w.RequestedTime = time.Now()

	return nil
}
