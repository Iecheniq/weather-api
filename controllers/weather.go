package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"
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
	jdata := models.WeatherJsonData{}
	weather := models.WeatherData{}
	city := o.GetString("city")
	country := o.GetString("country")
	if city == "" || country == "" {
		o.Data["json"] = errors.New("You must enter params city and country").Error()
	} else {
		cityCode := city + "," + country
		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%v&appid=1508a9a4840a5574c822d70ca2132032", cityCode)

		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			o.Data["json"] = err.Error()
		}
		req.Header.Add("accept-encoding", "application/json")
		response, err := client.Do(req)

		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Print(err)
				o.Data["json"] = err.Error()
			} else {
				json.Unmarshal(body, &jdata)
				if err := weather.GetDataFromJSON(jdata); err != nil {
					log.Print(err)
					o.Data["json"] = err.Error()

				} else {
					o.Data["json"] = weather
				}
			}
		}
	}

	o.ServeJSON()
}
