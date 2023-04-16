package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	"golang.org/x/exp/slog"

	"github.com/jkratz55/zapslog"
)

func main() {

	// Default slog Json logger
	sLogger := slog.New(slog.NewJSONHandler(os.Stderr))
	sLogger = sLogger.With(slog.String("appId", "FG2212"))
	sLogger.Info("Testing some stuffs",
		slog.Int("count", 11),
		slog.String("service", "test"),
		slog.String("err", "eeeeeee"),
		slog.Group("http",
			slog.Int("code", 123),
			slog.String("name", "cow")))

	zLogger, _ := zap.NewProduction()
	zSlogger := slog.New(zapslog.NewHandler(zLogger))
	zSlogger = zSlogger.With(slog.String("appId", "FG2212"))
	zSlogger.Info("Something something dark side")
	zSlogger.Log(context.Background(),
		zapslog.LevelError,
		"Testing some stuffs",
		slog.Int("count", 11),
		slog.String("service", "test"),
		slog.String("err", "eeeeeee"),
		slog.Any("data", map[string]any{
			"units": "cm",
		}),
		slog.Group("http",
			slog.Int("code", 123),
			slog.String("name", "cow")))

	zSlogger.Info("Did we panic?")

	// zSlogger.Log(context.Background(), slog_zap.LevelPanic, "Ohhhh snap",
	// 	slog.String("name", "something"),
	// 	slog.Group("test",
	// 		slog.String("transport", "http"),
	// 		slog.Group("http",
	// 			slog.String("method", "GET"))))
	//
	// sLogger.Log(context.Background(), slog_zap.LevelPanic, "Ohhhh snap",
	// 	slog.String("name", "something"),
	// 	slog.Group("test",
	// 		slog.String("transport", "http"),
	// 		slog.Group("http",
	// 			slog.String("method", "GET"))))
}
