package controllers

import (
	"fmt"
	"net/http"

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
	city := s.GetString("city")
	country := s.GetString("country")

	if city == "" || country == "" {
		http.Error(s.Ctx.ResponseWriter, "You must enter params 'city' and 'country' in the request body", http.StatusBadRequest)
		return
	}
	if err := models.AddScheduler(city, country); err != nil {
		if _, ok := err.(models.CityNotFoundError); ok {
			http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Ctx.ResponseWriter.WriteHeader(http.StatusAccepted)
	if _, err := s.Ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Scheduler for %v created", city))); err != nil {
		http.Error(s.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}
