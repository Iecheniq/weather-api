package models

import (
	"database/sql"
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
	Weather     []map[string]interface{} `json:"weather"`
	Main        map[string]float64       `json:"main"`
	Wind        map[string]float64       `json:"wind"`
	Sys         map[string]interface{}   `json:"sys"`
	City        string                   `json:"name"`
	Coordinates map[string]float64       `json:"coord"`
}

type WeatherData struct {
	Location      string
	Temperature   string
	Wind          map[string]float64
	Cloudines     string
	Presure       string
	Humidity      string
	Sunrise       time.Time
	Sunset        time.Time
	Coordinates   map[string]float64
	RequestedTime time.Time `json:"Requested_time"`
}

func SaveWeatherRequest(city, country string) error {
	stmt, err := db.Prepare("INSERT INTO weather_requests (city, country) VALUES (?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(city, country)
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

func IsRequestTimestampGreater(city string) (bool, error) {
	const requestTimeLimit float64 = 300
	requestTimestamp := ""
	err := db.QueryRow("SELECT time FROM weather_requests WHERE city = ? ORDER BY id DESC LIMIT 1", city).Scan(&requestTimestamp)
	if err != nil {
		return false, err
	}
	requestTime, err := time.Parse("2006-01-02 15:04:05", requestTimestamp)
	if err != nil {
		return false, err
	}
	if time.Since(requestTime).Seconds() > requestTimeLimit {
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
	wd.Location = op.City + ", " + op.Sys["country"].(string)
	wd.Temperature = strconv.FormatFloat(op.Main["temp"], 'f', 0, 64) + "°C"
	wd.Wind = op.Wind
	wd.Cloudines = op.Weather[0]["description"].(string)
	wd.Presure = strconv.FormatFloat(op.Main["pressure"], 'f', 0, 64) + " hpa"
	wd.Humidity = strconv.FormatFloat(op.Main["humidity"], 'f', 0, 64) + "%"
	wd.Sunrise = time.Unix(int64(op.Sys["sunrise"].(float64)), 0)
	wd.Sunset = time.Unix(int64(op.Sys["sunset"].(float64)), 0)
	wd.Coordinates = op.Coordinates
	wd.RequestedTime = time.Now()

	return wd, nil
}
