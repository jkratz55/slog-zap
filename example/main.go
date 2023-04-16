package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	"golang.org/x/exp/slog"

	slog_zap "github.com/jkratz55/slog-zap"
)

func main() {
	// opts := slog.HandlerOptions{ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
	// 	if attr.Key == slog.LevelKey {
	// 		level := attr.Value.Any().(slog.Level)
	// 		levelLabel := slog_zap.LevelName(level)
	//
	// 		attr.Value = slog.StringValue(levelLabel)
	// 	}
	//
	// 	return attr
	// }}

	sLogger := slog.New(slog.NewJSONHandler(os.Stderr))
	sLogger = sLogger.With(slog.String("appId", "FG2212"))
	sLogger = sLogger.WithGroup("main")
	sLogger.Info("Testing some stuffs",
		slog.Int("count", 11),
		slog.String("service", "test"),
		slog.String("err", "eeeeeee"),
		slog.Group("http",
			slog.Int("code", 123),
			slog.String("name", "cow")))

	zLogger, _ := zap.NewProduction()
	zSlogger := slog.New(slog_zap.NewHandler(zLogger))
	zSlogger.Log(context.Background(),
		slog_zap.LevelPanic,
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

	zLogger.Panic("")

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
