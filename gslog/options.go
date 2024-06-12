package gslog

type Options interface {
	apply(*LogHandlerOptions)
}

type optionFunc func(*LogHandlerOptions)

func (gs optionFunc) apply(options *LogHandlerOptions) {
	gs(options)
}

type LogHandlerOptions struct {
	// 日志级别信息
	level LevelEnabler

	// 文本格式化前缀
	textPrefix string
	// 格式化标记位
	textFlag int

	// json 时间格式化格式
	jsonTimeFormat string
}

func WithLevelEnabler(level LevelEnabler) Options {
	return optionFunc(func(options *LogHandlerOptions) {
		options.level = level
	})
}

func WithTextFlag(flag int) Options {
	return optionFunc(func(options *LogHandlerOptions) {
		options.textFlag = flag
	})
}
