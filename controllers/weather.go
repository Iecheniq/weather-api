package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/iecheniq/weather/external_services"
	"github.com/iecheniq/weather/models"
)

// Operations about city
type CityController struct {
	beego.Controller
}

// @Title Get
// @Description get city weather
// @Param city			query 	string	true		"City name "
// @Param country       query   string true "Country code"
// @Success 200 {string} models.City.Name
// @Failure 403 :City weather not found
// @router / [get]
func (o *CityController) Get() {
	weatherJData := models.WeatherJsonData{}
	weatherData := models.WeatherData{}
	city := o.GetString("city")
	country := o.GetString("country")
	if city == "" || country == "" {
		http.Error(o.Ctx.ResponseWriter, "You must enter params 'city' and 'country'", http.StatusBadRequest)
		return
	}
	response, err := services.GetWeather(city, country)
	if err != nil {
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Unmarshal(body, &weatherJData)
	if err := weatherData.GetDataFromJSON(weatherJData); err != nil {
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := models.SaveWeatherRequest(); err != nil {
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	o.Data["json"] = weatherData
	o.ServeJSON()
}
