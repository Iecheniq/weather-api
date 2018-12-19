package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/iecheniq/weather/models"
)

// Operations about weather
type WeatherController struct {
	beego.Controller
}

// @Title Get
// @Description get city weather
// @Param city			query 	string	true		"City name "
// @Param country       query   string true "Country code"
// @Success 200 {string} models.City.Name
// @Failure 403 :City weather not found
// @router / [get]
func (w *WeatherController) Get() {
	weatherJData := models.OpenWeatherJsonData{}
	city := w.GetString("city")
	country := w.GetString("country")
	if city == "" || country == "" {
		http.Error(w.Ctx.ResponseWriter, "You must enter params 'city' and 'country'", http.StatusBadRequest)
		return
	}
	response, err := models.GetWeather(city, country)
	if err != nil {
		http.Error(w.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		http.Error(w.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &weatherJData); err != nil {
		log.Print(err)
		http.Error(w.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	weatherData := weatherJData.ParseWeatherData()
	if err := models.SaveWeatherRequest(weatherData); err != nil {
		http.Error(w.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Data["json"] = *weatherData
	w.ServeJSON()
}
