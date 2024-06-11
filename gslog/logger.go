package gslog

import (
	"context"
	"runtime"
	"time"
)

type Logger struct {
	handler LogHandler
}

func NewLogger(handler LogHandler) *Logger {
	return &Logger{
		handler: handler,
	}
}

func (gs *Logger) Trace(msg string, args ...any) {
	gs.log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Debug(msg string, args ...any) {
	gs.log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Info(msg string, args ...any) {
	gs.log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Warn(msg string, args ...any) {
	gs.log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Error(msg string, args ...any) {
	gs.log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Critical(msg string, args ...any) {
	gs.log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	gs.log(ctx, TraceLevel, msg, args...)
}

func (gs *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	gs.log(ctx, DebugLevel, msg, args...)
}

func (gs *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	gs.log(ctx, InfoLevel, msg, args...)
}

func (gs *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	gs.log(ctx, WarnLevel, msg, args...)
}

func (gs *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	gs.log(ctx, ErrorLevel, msg, args...)
}

func (gs *Logger) CriticalContext(ctx context.Context, msg string, args ...any) {
	gs.log(ctx, CriticalLevel, msg, args...)
}

func (gs *Logger) TraceFields(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, TraceLevel, msg, args...)
}

func (gs *Logger) DebugFields(msg string, args ...Field) {
	gs.logFields(context.Background(), DebugLevel, msg, args...)
}

func (gs *Logger) InfoFields(msg string, args ...Field) {
	gs.logFields(context.Background(), InfoLevel, msg, args...)
}

func (gs *Logger) WarnFields(msg string, args ...Field) {
	gs.logFields(context.Background(), WarnLevel, msg, args...)
}

func (gs *Logger) ErrorFields(msg string, args ...Field) {
	gs.logFields(context.Background(), ErrorLevel, msg, args...)
}

func (gs *Logger) CriticalFields(msg string, args ...Field) {
	gs.logFields(context.Background(), CriticalLevel, msg, args...)
}

func (gs *Logger) TraceFieldsContext(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, TraceLevel, msg, args...)
}

func (gs *Logger) DebugFieldsContext(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, DebugLevel, msg, args...)
}

func (gs *Logger) InfoFieldsContext(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, InfoLevel, msg, args...)
}

func (gs *Logger) WarnFieldsContext(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, WarnLevel, msg, args...)
}

func (gs *Logger) ErrorFieldsContext(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, ErrorLevel, msg, args...)
}

func (gs *Logger) CriticalFieldsContext(ctx context.Context, msg string, args ...Field) {
	gs.logFields(ctx, CriticalLevel, msg, args...)
}

func (gs *Logger) Enable(ctx context.Context, level LevelEnabler) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	return gs.handler.Enabled(ctx, level)
}

func (gs *Logger) log(ctx context.Context, level LevelEnabler, msg string, args ...any) {
	if !gs.Enable(ctx, level) {
		return
	}

}

func (gs *Logger) logFields(ctx context.Context, level LevelEnabler, msg string, args ...Field) {
	if !gs.Enable(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// runtime.Callers. this function, this function's Caller
	runtime.Callers(3, pcs[:])

	entry := NewLogEntry(time.Now(), level.Level(), msg, pcs[0])
	entry.AppendFields(args...)
	if ctx == nil {
		ctx = context.Background()
	}

	_ = gs.handler.LogRecord(ctx, entry)
}
