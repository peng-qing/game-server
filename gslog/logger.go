package gslog

type Logger struct {
	handler LogHandler
}

func NewLogger(handler LogHandler) *Logger {
	return &Logger{
		handler: handler,
	}
}
