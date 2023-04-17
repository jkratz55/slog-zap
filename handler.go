package zapslog

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

type stacktrace struct{ pc [1]uintptr }

type fields struct{ fields []zapcore.Field }

// Handler is an implementation of slog.Handler interface backed by Zap's Logger.
//
// The zero-value of Handler is not usable. Handler should be created/initialized
// using the NewHandler function.
type Handler struct {
	logger     *zap.Logger
	stackPool  *sync.Pool
	fieldsPool *sync.Pool
}

// NewHandler creates and initializes a new Handler.
//
// This function will panic if passed a nil zap.Logger.
func NewHandler(logger *zap.Logger) *Handler {
	if logger == nil {
		panic("cannot create Handler with nil zap.Logger")
	}
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
	zapFields := make([]zapcore.Field, len(attrs))
	_ = h.fields(func(f []zapcore.Field) error {
		for _, attr := range attrs {
			f = h.appendAttr(f, attr, "")
		}
		copy(zapFields, f)
		return nil
	})

	return &Handler{
		logger:     h.logger.With(zapFields...),
		stackPool:  h.stackPool,
		fieldsPool: h.fieldsPool,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		logger:     h.logger.Named(name),
		stackPool:  h.stackPool,
		fieldsPool: h.fieldsPool,
	}
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
