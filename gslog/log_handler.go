package gslog

import "context"

type LogHandler interface {
	// Enabled 针对每个Handler支持处理不同的日志级别
	Enabled(ctx context.Context, levelEnabler LevelEnabler) bool
	// LogRecord 处理每条日志元数据
	LogRecord(ctx context.Context, entry *LogEntry) error
}
