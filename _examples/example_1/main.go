package main

import (
	"fmt"
	"github.com/fnsne/irisslog"
	"github.com/kataras/iris/v12"
	"golang.org/x/exp/slog"
	"os"
	"time"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	slog.SetDefault(logger)

	app := iris.New()
	app.Use(irisslog.New(logger))

	// Example ping request.
	app.Get("/ping", func(ctx iris.Context) {
		handlerLogger := irisslog.GetLogger(ctx)
		handlerLogger.Info("info ping")
		handlerLogger.Debug("debug ping")
		ctx.Text("pong " + fmt.Sprint(time.Now().Unix()))
	})

	// Example when panic happen.
	app.Get("/panic", func(ctx iris.Context) {
		panic("An unexpected error happen!")
	})

	// Listen and Server in 0.0.0.0:8080
	app.Listen(":8081")

}
