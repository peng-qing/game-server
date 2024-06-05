package gslog

// Option options模式
type Option interface {
	apply(logger *Logger)
}

type optionFunc func(logger *Logger)

func (gs optionFunc) apply(logger *Logger) {
	gs(logger)
}

func WithFlags(flag int) Option {
	return optionFunc(func(logger *Logger) {
		logger.SetFlags(flag)
	})
}

func WithPrefix(prefix string) Option {
	return optionFunc(func(logger *Logger) {
		logger.SetPrefix(prefix)
	})
}

func WithLevelEnabler(enabler LevelEnabler) Option {
	return optionFunc(func(logger *Logger) {
		logger.SetLevelEnabler(enabler)
	})
}

func WithWriteSyncer(name string, syncer WriteSyncer) Option {
	return optionFunc(func(logger *Logger) {
		logger.AppendWriteSyncer(name, syncer)
	})
}
