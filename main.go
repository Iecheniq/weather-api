package main

import (
	"log"

	_ "github.com/iecheniq/weather/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/weather_db?charset=utf8") //change localhost to weather
}

func main() {

	name := "default"
	force := false
	verbose := false
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		log.Fatal(err)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
