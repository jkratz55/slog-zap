package main

import (
	"go.mrchanchal.com/zaphandler"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
)

func main() {
	zapL, _ := zap.NewProduction()
	defer zapL.Sync()

	logger := slog.New(zaphandler.New(zapL))

	logger.Info("Testing some stuffs",
		slog.Int("count", 11),
		slog.String("service", "test"),
		slog.String("err", "eeeeeee"),
		slog.Group("http",
			slog.Int("code", 123),
			slog.String("name", "cow")))

}
