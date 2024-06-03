package gslog

// Option options模式
type Option interface {
	apply(logger *Logger)
}

type optionFunc func(logger *Logger)

func (gs optionFunc) apply(logger *Logger) {
	gs(logger)
}
