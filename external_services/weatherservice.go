package services

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
)

func GetWeather(city, country string) (*http.Response, error) {

	cityCode := city + "," + country
	url := fmt.Sprintf(beego.AppConfig.String("ExternalAPIWeatherURL"), cityCode)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
