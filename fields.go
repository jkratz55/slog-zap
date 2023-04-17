package zapslog

import (
	"fmt"
	"math"

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

func getType(val slog.Value) (zapcore.FieldType, int64, string, any) {
	var (
		Type      = zapcore.ReflectType
		Integer   = int64(0)
		String    = ""
		Interface = any(nil)
	)

	Kind := val.Kind()
	if t, found := mapField(Kind); found {
		Type = t
	}

	switch Kind {
	case slog.KindAny:
		Interface = val.Any()

		switch Interface.(type) {
		case fmt.Stringer:
			Type = zapcore.StringerType
		case error:
			Type = zapcore.ErrorType
		}
	case slog.KindBool:
		if val.Bool() {
			Integer = 1
		}
	case slog.KindDuration:
		Integer = int64(val.Duration())
	case slog.KindFloat64:
		Integer = int64(math.Float64bits(val.Float64()))
	case slog.KindInt64:
		Integer = val.Int64()
	case slog.KindString:
		String = val.String()
	case slog.KindTime:
		t := val.Time()
		Integer = t.UnixNano()
		Interface = t.Location()
	case slog.KindUint64:
		Type = zapcore.Uint64Type
		Integer = int64(val.Uint64())
	case slog.KindGroup:
	case slog.KindLogValuer:
	}

	return Type, Integer, String, Interface
}
