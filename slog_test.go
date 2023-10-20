package irisslog_test

import (
	irisslog "github.com/iris-contrib/middleware/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kataras/iris/v12"
)

func TestNew(t *testing.T) {
	//todo 要用slog的話，要mock handler，來檢查是否有發送log出來。
	app := iris.New()

	utcLogger, utcLoggerObserved := buildDummyLogger()
	app.Use(irisslog.New(utcLogger, time.RFC3339, true))

	localLogger, localLoggerObserved := buildDummyLogger()
	app.Use(irisslog.New(localLogger, time.RFC3339, false))

	app.Get("/test", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusNoContent)
		ctx.JSON(nil)
	})

	if err := app.Build(); err != nil {
		t.Fatal(err)
	}

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	app.ServeHTTP(res1, req1)

	//todo check log one log

	//todo check log.path right

}
