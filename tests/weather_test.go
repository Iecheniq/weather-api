package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/iecheniq/weather/models"
	_ "github.com/iecheniq/weather/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestWeatherGet(t *testing.T) {
	db := models.MySQLWeatherDb{
		DataSource: "root:root@tcp(localhost:3306)/weather_db_test",
	}

	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	testCases := []struct {
		name string
		url  string
		code int
	}{
		{name: "Get weather with correct params", url: "/weather?city=Mexico&country=mx", code: 200},
		{name: "Get weather with missing params", url: "/weather", code: 400},
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
