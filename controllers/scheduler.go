package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	var response *http.Response
	var err error
	notFound := false
	done := make(chan bool)
	if city == "" || country == "" {
		http.Error(s.Ctx.ResponseWriter, "You must enter params 'city' and 'country' in the request body", http.StatusBadRequest)
		return
	}
	go func() {
		for {
			response, err = models.GetWeather(city, country)
			if response.StatusCode == http.StatusNotFound {
				notFound = true
				break
			}
			if err != nil {
				break
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				break
			}
			if err = json.Unmarshal(body, &weatherJData); err != nil {
				break
			}
			weatherData, err := weatherJData.ParseWeatherData()
			if err != nil {
				break
			}
			if err = models.SaveWeatherRequest(weatherData); err != nil {
				break
			}
			done <- true
			time.Sleep(60 * time.Second)
		}
		done <- true
	}()
	<-done
	if notFound {
		fmt.Printf("City not found")
		http.Error(s.Ctx.ResponseWriter, errors.New("Cityt not found").Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Print(err)
		http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Ctx.ResponseWriter.WriteHeader(http.StatusAccepted)
	if _, err := s.Ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Scheduler for %v created", city))); err != nil {
		http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}
