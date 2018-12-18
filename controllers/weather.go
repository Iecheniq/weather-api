package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"
	services "github.com/iecheniq/weather/external_services"
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
	weatherJData := models.OpenWeatherJsonData{}
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
	if err := json.Unmarshal(body, &weatherJData); err != nil {
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	weatherData, err := weatherJData.ParseWeatherData()
	if err != nil {
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	saveRequest, err := models.IsRequestTimestampGreater(city)
	if err != nil {
		if err == sql.ErrNoRows {
			saveRequest = true
			goto save_request
		}
		log.Print(err)
		http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
save_request:
	if saveRequest {
		if err := models.SaveWeatherRequest(city, country); err != nil {
			log.Print(err)
			http.Error(o.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	o.Data["json"] = *weatherData
	o.ServeJSON()
}
