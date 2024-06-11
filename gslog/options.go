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
	// 格式化标记位
	flag int

	// 文本格式化前缀
	textPrefix string
}
