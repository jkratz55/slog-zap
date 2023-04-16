package zapslog

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

type stacktrace struct{ pc [1]uintptr }

type fields struct{ fields []zapcore.Field }

type Handler struct {
	logger     *zap.Logger
	stackPool  *sync.Pool
	fieldsPool *sync.Pool
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger:     logger,
		stackPool:  &sync.Pool{New: func() any { return &stacktrace{pc: [1]uintptr{}} }},
		fieldsPool: &sync.Pool{New: func() any { return &fields{fields: make([]zapcore.Field, 0, 32)} }},
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.Core().Enabled(mapLevel(level))
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	frame, _ := h.frame(record.PC)
	lvl := mapLevel(record.Level)

	entry := zapcore.Entry{
		Level:   lvl,
		Time:    record.Time,
		Message: record.Message,
		Caller:  zapcore.NewEntryCaller(frame.PC, frame.File, frame.Line, true),
		Stack:   "",
	}

	checked := h.logger.Core().Check(entry, nil)
	if checked == nil {
		return nil
	}

	// If log level is panic call panic after writing the entry
	if lvl == zapcore.PanicLevel {
		checked.After(entry, zapcore.WriteThenPanic)
	}

	return h.fields(func(f []zapcore.Field) error {
		record.Attrs(func(attr slog.Attr) {
			f = h.appendAttr(f, attr, "")
		})
		checked.Write(f...)
		return nil
	})
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) WithGroup(name string) slog.Handler {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) frame(caller uintptr) (runtime.Frame, error) {
	stack, ok := h.stackPool.Get().(*stacktrace)
	if !ok {
		return runtime.Frame{}, fmt.Errorf("failed to retrieve stack from pool")
	}
	defer h.stackPool.Put(stack)
	stack.pc[0] = caller
	frames, _ := runtime.CallersFrames(stack.pc[:]).Next()
	return frames, nil
}

func (h *Handler) fields(fn func([]zapcore.Field) error) error {
	fields, ok := h.fieldsPool.Get().(*fields)
	if !ok {
		return fmt.Errorf("failed to retrieve fields from pool")
	}
	defer h.fieldsPool.Put(fields)
	return fn(fields.fields[:0])
}

func (h *Handler) appendAttr(fields []zapcore.Field, attr slog.Attr, prefix string) []zapcore.Field {
	if attr.Value.Kind() != slog.KindGroup {
		Type, Integer, String, Interface := getType(attr.Value)
		return append(fields, zapcore.Field{
			Key:       prefix + attr.Key,
			Type:      Type,
			Integer:   Integer,
			String:    String,
			Interface: Interface,
		})
	}

	for _, gAttr := range attr.Value.Group() {
		fields = h.appendAttr(fields, gAttr, attr.Key+".")
	}

	return fields
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
