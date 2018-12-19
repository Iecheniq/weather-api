package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/iecheniq/weather/models"
)

// Operations about weather
type SchedulerController struct {
	beego.Controller
}

// @Title PUT
// @Description Add a new shceduler to get the weather of a city every one hour
// @Param city			body 	string	true		"City name "
// @Param country       body   string true "Country code"
// @Success 202 {string} models.City.Name
// @Failure 403 :City weather not found
// @router / [put]
func (s *SchedulerController) Put() {
	weatherJData := models.OpenWeatherJsonData{}
	city := s.GetString("city")
	country := s.GetString("country")
	if city == "" || country == "" {
		http.Error(s.Ctx.ResponseWriter, "You must enter params 'city' and 'country' in the request body", http.StatusBadRequest)
		return
	}
	go func() {
		for {
			response, err := models.GetWeather(city, country)
			if err != nil {
				http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Print(err)
				http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := json.Unmarshal(body, &weatherJData); err != nil {
				log.Print(err)
				http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
			weatherData := weatherJData.ParseWeatherData()
			if err := models.SaveWeatherRequest(weatherData); err != nil {
				http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
			time.Sleep(60 * time.Second)
		}
	}()
	s.Ctx.ResponseWriter.WriteHeader(http.StatusAccepted)
	if _, err := s.Ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Scheduler for %v created", city))); err != nil {
		http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

}
