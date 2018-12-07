package models

import (
	"errors"
	"strings"
)

var (
	Cities map[string]*City
)

type City struct {
	Name    string
	Country string
}

func init() {
	Cities = make(map[string]*City)
	Cities["Bogota"] = &City{Name: "Bogota", Country: "co"}
	Cities["Mexico"] = &City{Name: "Mexico", Country: "mx"}
}

func GetCity(name string) (*City, error) {
	name = strings.ToLower(name)
	name = strings.Title(name)
	if city, ok := Cities[name]; ok {
		return city, nil
	}
	return nil, errors.New("City Does Not Exist")

}

func GetCityCode(city *City) string {
	return city.Name + "," + city.Country
}
