package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type WeatherDb interface {
	Open()
	Close()
}

type WeatherJSONData interface {
	ParseWeatherData() (*WeatherData, error)
}

type MySQLWeatherDb struct {
	DataSource string
}

type OpenWeatherJsonData struct {
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

func SaveWeatherRequest() error {
	res, err := db.Exec("INSERT INTO weather_requests () VALUES ()")
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Affected = %d\n", rowCnt)
	return nil
}

func IsRequestTimestampGreater() (bool, error) {
	timestamp := time.Time{}
	err := db.QueryRow("SELECT * FROM weather_requests ORDER BY id DESC LIMIT 1").Scan(&timestamp)
	if err != nil {
		return false, err
	}
	if time.Now().Second()-timestamp.Second() > 300 {
		return true, nil
	}
	return false, nil
}

func (database *MySQLWeatherDb) Open() error {
	db, _ = sql.Open("mysql", database.DataSource)

	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (database *MySQLWeatherDb) Close() {
	db.Close()
}

func (op *OpenWeatherJsonData) ParseWeatherData() (wd *WeatherData, err error) {

	wd = &WeatherData{}
	country, ok := op.Sys["country"].(string)
	if !ok {
		err = fmt.Errorf("Country is not type string")
		return nil, err
	}
	wd.Location = op.City + ", " + country
	wd.Temperature = strconv.FormatFloat(op.Main["temp"], 'f', 0, 64) + "°C"
	wd.Wind = op.Wind
	wd.Cloudines = op.Weather[0]["description"]
	wd.Presure = strconv.FormatFloat(op.Main["pressure"], 'f', 0, 64) + " hpa"
	wd.Humidity = strconv.FormatFloat(op.Main["humidity"], 'f', 0, 64) + "%"
	wd.Sunrise = op.Sys["sunrise"].(float64)
	wd.Sunset = op.Sys["sunset"].(float64)
	wd.Coordinates = op.Coordinates
	wd.RequestedTime = time.Now()

	return wd, nil
}
