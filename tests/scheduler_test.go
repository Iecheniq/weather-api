package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/iecheniq/weather/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/weather_db_test?charset=utf8") //change localhost to weather
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
func TestSchedulerPut(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		body map[string]string
		code int
	}{
		{name: "Add scheduler with correct params in body", url: "/weather/scheduler", body: map[string]string{"city": "Mexico", "country": "mx"}, code: 202},
		{name: "Add Scheduler with missing params in body", url: "/weather/scheduler", body: map[string]string{"city": "Mexico"}, code: 400},
		{name: "Add Scheduler with wrong params in body", url: "/weather/scheduler", body: map[string]string{"city": "Mico", "country": "ml"}, code: 404},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.body)
			if err != nil {
				log.Fatal(err)
			}
			r, err := http.NewRequest("PUT", tc.url, bytes.NewBuffer(payload))
			if err != nil {
				fmt.Print(err)
			}
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)
			beego.Trace("testing", "TestSchedulerPut", "Code[%d]\n%s", w.Code, w.Body.String())

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
