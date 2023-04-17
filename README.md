# zapslog

ZapSlog provides a `slog.Handler` implementation backed by Uber's `zap.Logger`.

_Notice: `slog` is still experimental and not officially part of the standard library. There may be changes to `slog` before it becomes finalized and part of the standard library. You may want to keep that in mind before using `slog` or this library._

## Why slog?

`slog` provides structured logging functionality from the Go team. Although experimental at this time, ideally it will become the defacto logging interface used by frameworks and libraries. Presently the logging ecosystem in Go is fragmented although there are a few excellent third party libraries such as Zap and Zerolog. Third party loggers can be used as Handlers for `slog` making it very versatile for library developers.

## Why zapslog?

`zapslog` uses `zap.Logger` as the handler for slog to get the performance Zap is known for while using the `slog` interfaces.

## Example

```go
package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	"golang.org/x/exp/slog"

	"github.com/jkratz55/zapslog"
)

func main() {

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
}
```