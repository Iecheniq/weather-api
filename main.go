package main

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/iecheniq/weather/models"
	_ "github.com/iecheniq/weather/routers"
)

func main() {
	db := models.MySQLWeatherDb{
		DataSource: "root:root@tcp(localhost:3306)/weather_db",
	}

	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
