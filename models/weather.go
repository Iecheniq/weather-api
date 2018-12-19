package models

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	services "weather/external_services"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type WeatherJSONData interface {
	ParseWeatherData() (*WeatherData, error)
}

type OpenWeatherJsonData struct {
	Weather     []map[string]interface{} `json:"weather"`
	Main        map[string]float64       `json:"main"`
	Wind        map[string]float64       `json:"wind"`
	Sys         map[string]interface{}   `json:"sys"`
	City        string                   `json:"name"`
	Coordinates map[string]float64       `json:"coord"`
}

type WeatherData struct {
	Id            int64 `json:"-"`
	Location      string
	Temperature   string
	Wind          string
	Cloudines     string
	Presure       string
	Humidity      string
	Sunrise       time.Time
	Sunset        time.Time
	Coordinates   string
	RequestedTime time.Time `json:"Requested_time"`
}

func init() {
	orm.RegisterModel(new(WeatherData))
}

func GetWeather(city, country string) (*http.Response, error) {
	//TODO Support different weather services
	response, err := services.GetWeather(city, country)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return response, nil
}

func saveWeatherRequestDb(wd *WeatherData) error {
	o := orm.NewOrm()
	id, err := o.Insert(wd)
	if err != nil {
		fmt.Print(err)
		return err
	}
	fmt.Printf("Weather data with ID %v created", id)
	return nil
}

func isRequestTimestampGreater(location string) (bool, error) {

	const requestTimeLimit float64 = 300
	wd := WeatherData{}
	o := orm.NewOrm()
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return false, err
	}
	qb.Select("id, requested_time").
		From("weather_data").
		Where("location = ?").
		OrderBy("id").
		Desc().
		Limit(1)
	query := qb.String()
	err = o.Raw(query, location).QueryRow(&wd)
	if err != nil {
		return false, err
	}
	if time.Since(wd.RequestedTime).Seconds() > requestTimeLimit {
		return true, nil
	}
	return false, nil
}

func SaveWeatherRequest(w *WeatherData) error {
	saveRequest, err := isRequestTimestampGreater(w.Location)
	if err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			saveRequest = true
			goto save_request
		}
		log.Print(err)
		return err
	}
save_request:
	if saveRequest {
		if err := saveWeatherRequestDb(w); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func (op *OpenWeatherJsonData) ParseWeatherData() (*WeatherData, error) {

	wd := &WeatherData{}
	country, ok := op.Sys["country"].(string)
	if !ok {
		err := fmt.Errorf("Country is not type string")
		return nil, err
	}
	cloudines, ok := op.Weather[0]["description"].(string)
	if !ok {
		err := fmt.Errorf("Cloudines is not type string")
		return nil, err
	}
	sunrise, ok := op.Sys["sunrise"].(float64)
	if !ok {
		err := fmt.Errorf("Sunrise is not type float64")
		return nil, err
	}
	sunset, ok := op.Sys["sunset"].(float64)
	if !ok {
		err := fmt.Errorf("Sunset is not type float64")
		return nil, err
	}
	wd.Location = op.City + ", " + country
	wd.Temperature = strconv.FormatFloat(op.Main["temp"]-273.15, 'f', 0, 64) + "°C"
	wd.Wind = strconv.FormatFloat(op.Wind["speed"], 'f', 0, 64) + "m/s, " + strconv.FormatFloat(op.Wind["deg"], 'f', 0, 64) + "°"
	wd.Cloudines = cloudines
	wd.Presure = strconv.FormatFloat(op.Main["pressure"], 'f', 0, 64) + " hpa"
	wd.Humidity = strconv.FormatFloat(op.Main["humidity"], 'f', 0, 64) + "%"
	wd.Sunrise = time.Unix(int64(sunrise), 0)
	wd.Sunset = time.Unix(int64(sunset), 0)
	wd.Coordinates = strconv.FormatFloat(op.Coordinates["lon"], 'f', 0, 64) + ", " + strconv.FormatFloat(op.Coordinates["lat"], 'f', 0, 64)
	wd.RequestedTime = time.Now()

	return wd, nil
}
