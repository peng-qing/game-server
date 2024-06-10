package gslog

import "context"

type Logger struct {
	handler LogHandler
}

func NewLogger(handler LogHandler) *Logger {
	return &Logger{
		handler: handler,
	}
}

func (gs *Logger) Trace(msg string, args ...any) {
	gs.Log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Debug(msg string, args ...any) {
	gs.Log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Info(msg string, args ...any) {
	gs.Log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Warn(msg string, args ...any) {
	gs.Log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Error(msg string, args ...any) {
	gs.Log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) Critical(msg string, args ...any) {
	gs.Log(context.Background(), TraceLevel, msg, args...)
}

func (gs *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	gs.Log(ctx, TraceLevel, msg, args...)
}

func (gs *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	gs.Log(ctx, DebugLevel, msg, args...)
}

func (gs *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	gs.Log(ctx, InfoLevel, msg, args...)
}

func (gs *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	gs.Log(ctx, WarnLevel, msg, args...)
}

func (gs *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	gs.Log(ctx, ErrorLevel, msg, args...)
}

func (gs *Logger) CriticalContext(ctx context.Context, msg string, args ...any) {
	gs.Log(ctx, CriticalLevel, msg, args...)
}

func (gs *Logger) Log(ctx context.Context, level LogLevel, msg string, args ...any) {

}
