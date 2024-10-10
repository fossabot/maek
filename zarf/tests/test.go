package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"time"

	"github.com/bluele/go-timecop"

	"github.com/beego/beego/v2/client/orm"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/karngyan/maek/conf"
	"github.com/karngyan/maek/db"
	"github.com/karngyan/maek/domains"
	"github.com/karngyan/maek/routers"
)

var frozenTime = time.Unix(1234567890, 0)
var initOnce sync.Once

func InitApp() error {
	initFn := func() error {
		log := logs.NewLogger()
		defer log.Flush()

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root := os.Getenv("ROOT")
		web.TestBeegoInit(root)

		// beego changes dir internally after TestBeegoInit
		// we'd like it back for approvals to work
		err = os.Chdir(cwd)
		if err != nil {
			return err
		}

		if err := routers.Init(log); err != nil {
			return err
		}

		if err := conf.Init(); err != nil {
			return err
		}

		if err := db.Init(); err != nil {
			panic(err)
		}

		if err := domains.InitTest(); err != nil {
			return err
		}

		return nil
	}

	var initErr error

	initOnce.Do(func() {
		initErr = initFn()
	})

	return initErr
}

func FreezeTime() {
	timecop.Freeze(frozenTime)
}

func CleanDBRows() {
	// force would drop the tables and recreate them
	_ = orm.RunSyncdb("default", true, false)
}

func Post(path string, body any) (*httptest.ResponseRecorder, error) {
	return request(http.MethodPost, path, body)
}

func Get(path string) (*httptest.ResponseRecorder, error) {
	return request(http.MethodGet, path, nil)
}

func Put(path string, body any) (*httptest.ResponseRecorder, error) {
	return request(http.MethodPut, path, body)
}

func Patch(path string, body any) (*httptest.ResponseRecorder, error) {
	return request(http.MethodPatch, path, body)
}

func Delete(path string) (*httptest.ResponseRecorder, error) {
	return request(http.MethodDelete, path, nil)
}

func request(method string, path string, body any) (*httptest.ResponseRecorder, error) {
	buf := bytes.NewBuffer([]byte{})
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	req, err := http.NewRequest(method, path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	web.BeeApp.Handlers.ServeHTTP(rr, req)

	return rr, nil
}