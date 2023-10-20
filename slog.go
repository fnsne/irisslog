// Package iriszap provides log handling using zap package.
package irisslog

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/kataras/iris/v12"
	"golang.org/x/exp/slog"
	"time"
)

func New(logger *slog.Logger) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		start := time.Now()
		path := ctx.Request().URL.Path
		query := ctx.Request().URL.RawQuery

		traceID, err := uuid.NewV4()
		if err != nil {
			logger.Warn("generate trace id failed", slog.String("error", err.Error()))
		}
		ctx.Values().Set("traceID", traceID.String())
		//to use traceID in handler, GetLogger(ctx)
		ctx.Next()

		// after request
		end := time.Now()
		latency := end.Sub(start)
		if ctx.GetErr() != nil {
			logger.Error(
				ctx.GetErr().Error(),
				slog.Int("status", ctx.GetStatusCode()),
				slog.String("method", ctx.Method()),
				slog.String("path", path),
				slog.String("query", query),
				slog.Duration("latency", latency),
			)
		} else {
			logger.Info(
				"Success",
				slog.Int("status", ctx.GetStatusCode()),
				slog.String("method", ctx.Method()),
				slog.String("path", path),
				slog.String("query", query),
				slog.Duration("latency", latency),
			)
		}
	}
}

func GetLogger(ctx iris.Context) *slog.Logger {
	logger := slog.With("trace_id", ctx.Values().GetString("traceID"))
	return logger
}
