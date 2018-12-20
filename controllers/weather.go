package controllers

import (
	"net/http"

	"github.com/iecheniq/weather/models"

	"github.com/astaxie/beego"
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
	city := w.GetString("city")
	country := w.GetString("country")
	if city == "" || country == "" {
		http.Error(w.Ctx.ResponseWriter, "You must enter params 'city' and 'country'", http.StatusBadRequest)
		return
	}
	weatherData, err := models.GetWeather(city, country)
	if err != nil {
		if _, ok := err.(models.CityNotFoundError); ok {
			http.Error(w.Ctx.ResponseWriter, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Data["json"] = *weatherData
	w.ServeJSON()
}
