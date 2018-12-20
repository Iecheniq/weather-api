package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var service WeatherService

type WeatherService interface {
	GetData(city, country string) error
	ParseData() (*WeatherData, error)
}

type CityNotFoundError struct {
	Description string
}

type OpenWeatherService struct {
	URL         string
	Weather     []map[string]interface{} `json:"weather"`
	Main        map[string]float64       `json:"main"`
	Wind        map[string]float64       `json:"wind"`
	Sys         map[string]interface{}   `json:"sys"`
	City        string                   `json:"name"`
	Coordinates map[string]float64       `json:"coord"`
	mux         sync.Mutex
}

type WeatherFromFilesService struct {
	Route       string
	Location    string `json:"location"`
	Temperature string `json:"temp"`
	Wind        string `json:"wind"`
	Cloudines   string `json:"cloudines"`
	Pressure    string `json:"pressure"`
	Humidity    string `json:"humidity"`
	Sunrise     int64  `json:"sunrise"`
	Sunset      int64  `json:"sunset"`
	Coordinates string `json:"coord"`
	mux         sync.Mutex
}

type WeatherData struct {
	Id            int64 `json:"-"`
	Location      string
	Temperature   string
	Wind          string
	Cloudines     string
	Pressure      string
	Humidity      string
	Sunrise       time.Time
	Sunset        time.Time
	Coordinates   string
	RequestedTime time.Time `json:"Requested_time"`
}

func init() {
	orm.RegisterModel(new(WeatherData))
	if err := setWeatherService(); err != nil {
		log.Fatal(err)
	}
}

func setWeatherService() error {
	selectedService := beego.AppConfig.String("WeatherService")
	switch selectedService {
	case "open_weather":
		service = &OpenWeatherService{URL: beego.AppConfig.String("OpenWeatherURL")}
	case "from_files":
		service = &WeatherFromFilesService{}
	default:
		err := errors.New("Configuration variable WeatherService not set or not supported.\nValue options: open_weather, from_files")
		fmt.Print(err)
		return err
	}
	return nil
}

func GetWeather(city, country string) (*WeatherData, error) {
	workers := make(chan error, 5)
	go func() {
		workers <- service.GetData(city, country)
	}()
	err := <-workers
	if err != nil {
		return nil, err
	}
	wd, err := service.ParseData()
	if err != nil {
		return nil, err
	}
	if err := SaveWeatherRequest(wd); err != nil {
		return nil, err
	}
	return wd, nil
}

func AddScheduler(city, country string) error {

	done := make(chan error)
	go func() {
		var err error
		for {
			_, err = GetWeather(city, country)
			if err != nil {
				break
			}
			done <- err
			time.Sleep(3600 * time.Second)
		}
		done <- err
	}()
	return <-done
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

func (e CityNotFoundError) Error() string {
	return e.Description
}

func (ow *OpenWeatherService) GetData(city, country string) error {
	ow.mux.Lock()
	defer ow.mux.Unlock()
	cityCode := city + "," + country
	url := fmt.Sprintf(ow.URL, cityCode)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		err = CityNotFoundError{Description: "City not found"}
		log.Print(err)
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return err
	}
	if err := json.Unmarshal(body, ow); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (op OpenWeatherService) ParseData() (*WeatherData, error) {

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
	wd.Pressure = strconv.FormatFloat(op.Main["pressure"], 'f', 0, 64) + " hpa"
	wd.Humidity = strconv.FormatFloat(op.Main["humidity"], 'f', 0, 64) + "%"
	wd.Sunrise = time.Unix(int64(sunrise), 0)
	wd.Sunset = time.Unix(int64(sunset), 0)
	wd.Coordinates = strconv.FormatFloat(op.Coordinates["lon"], 'f', 0, 64) + ", " + strconv.FormatFloat(op.Coordinates["lat"], 'f', 0, 64)
	wd.RequestedTime = time.Now()

	return wd, nil
}

func (wf *WeatherFromFilesService) GetData(city, country string) error {
	wf.mux.Lock()
	defer wf.mux.Unlock()
	filename := "w_" + city + "_" + country + ".json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, wf); err != nil {
		return err
	}
	return nil
}
func (wf WeatherFromFilesService) ParseData() (*WeatherData, error) {

	wd := &WeatherData{}
	wd.Location = wf.Location
	wd.Temperature = wf.Temperature
	wd.Wind = wf.Wind
	wd.Cloudines = wf.Cloudines
	wd.Pressure = wf.Pressure
	wd.Humidity = wf.Humidity
	wd.Sunrise = time.Unix(wf.Sunrise, 0)
	wd.Sunset = time.Unix(wf.Sunset, 0)
	wd.Coordinates = wf.Coordinates
	wd.RequestedTime = time.Now()

	return wd, nil
}
