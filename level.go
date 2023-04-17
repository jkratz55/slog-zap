package zapslog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

const (
	LevelDebug  = slog.LevelDebug
	LevelInfo   = slog.LevelInfo
	LevelWarn   = slog.LevelWarn
	LevelError  = slog.LevelError
	LevelDPanic = slog.Level(9)
	LevelPanic  = slog.Level(10)
	LevelFatal  = slog.Level(11)
)

var levelNames = map[slog.Leveler]string{
	LevelDebug:  "DEBUG",
	LevelInfo:   "INFO",
	LevelWarn:   "WARN",
	LevelError:  "ERROR",
	LevelDPanic: "DPANIC",
	LevelPanic:  "PANIC",
	LevelFatal:  "FATAL",
}

func LevelName(leveler slog.Leveler) string {
	return levelNames[leveler]
}

func mapLevel(lvl slog.Level) zapcore.Level {
	switch lvl {
	case slog.LevelDebug:
		return zap.DebugLevel
	case slog.LevelInfo:
		return zap.InfoLevel
	case slog.LevelWarn:
		return zap.WarnLevel
	case slog.LevelError:
		return zap.ErrorLevel
	case LevelPanic:
		return zap.PanicLevel
	case LevelFatal:
		return zap.FatalLevel
	default:
		// If there is no mapping from slog level to zap level default to error
		// level and hope someone reviewing logs notices
		return zap.ErrorLevel
	}
}
