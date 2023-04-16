package zapslog

import (
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

var slogToZapTypeMap = map[slog.Kind]zapcore.FieldType{
	slog.KindBool:     zapcore.BoolType,
	slog.KindDuration: zapcore.DurationType,
	slog.KindFloat64:  zapcore.Float64Type,
	slog.KindInt64:    zapcore.Int64Type,
	slog.KindString:   zapcore.StringType,
	slog.KindTime:     zapcore.TimeType,
	slog.KindUint64:   zapcore.Uint64Type,
}

func mapField(kind slog.Kind) (zapcore.FieldType, bool) {
	val, ok := slogToZapTypeMap[kind]
	return val, ok
}
