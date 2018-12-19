package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "weather/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(weather_db:3306)/weather_db_test?charset=utf8")
	name := "default"
	force := false
	verbose := false
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		log.Fatal(err)
	}
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestWeatherGet(t *testing.T) {

	testCases := []struct {
		name string
		url  string
		code int
	}{
		{name: "Get weather with correct params", url: "/weather?city=Mexico&country=mx", code: 200},
		{name: "Get weather with missing params", url: "/weather", code: 400},
		{name: "Get weather with wrong params", url: "/weather?city=Mexo&country=mc", code: 404},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)
			beego.Trace("testing", "TestWeatherGet", "Code[%d]\n%s", w.Code, w.Body.String())

			Convey("Subject: Test Station Endpoint\n", t, func() {
				Convey(fmt.Sprintf("Status Code Should Be %v", tc.code), func() {
					So(w.Code, ShouldEqual, tc.code)
				})
				Convey("The Result Should Not Be Empty", func() {
					So(w.Body.Len(), ShouldBeGreaterThan, 0)
				})
			})
		})
	}
}
