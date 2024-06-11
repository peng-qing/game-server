package gslog

import (
	"bytes"
	"errors"
	"fmt"
)

// LogLevel 日志级别有了部分调整
// 1. 默认TraceLevel 为0，保持默认级别为int默认值

var (
	errUnmarshalInvalid = errors.New("unmarshal invalid text")
)

type LogLevel int

const (
	TraceLevel LogLevel = iota << 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

// ParseLogLevel 解析字符串到LogLevel
func ParseLogLevel(text string) (LogLevel, error) {
	var lv LogLevel
	err := lv.UnmarshalText([]byte(text))

	return lv, err
}

// UnmarshalText 解析文本为指定日志级别
func (gs LogLevel) UnmarshalText(text []byte) error {
	if !gs.unmarshalText(text) && !gs.unmarshalText(bytes.ToUpper(text)) {
		return errUnmarshalInvalid
	}
	return nil
}

func (gs *LogLevel) unmarshalText(text []byte) bool {
	switch string(text) {
	case "trace", "TRACE":
		*gs = TraceLevel
	case "debug", "DEBUG":
		*gs = DebugLevel
	case "info", "INFO":
		*gs = InfoLevel
	case "warn", "WARN":
		*gs = WarnLevel
	case "error", "ERROR":
		*gs = ErrorLevel
	case "critical", "CRITICAL":
		*gs = CriticalLevel
	default:
		return false
	}
	return true
}

// LowCaseString 小写字母形式
// 也可以采用数组的方式 从0开始对应各个级别的下标字符串
func (gs LogLevel) LowCaseString() string {
	switch gs {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case CriticalLevel:
		return "critical"
	default:
		return fmt.Sprintf("LogLevel({%d})", gs)
	}
}

// UpCaseString 大写字母形式
func (gs LogLevel) UpCaseString() string {
	switch gs {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	default:
		return fmt.Sprintf("LogLevel({%d})", gs)
	}
}

// CapitalString 首字母大写形式
func (gs LogLevel) CapitalString() string {
	switch gs {
	case TraceLevel:
		return "Trace"
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	case ErrorLevel:
		return "Error"
	case CriticalLevel:
		return "Critical"
	default:
		return fmt.Sprintf("LogLevel({%d})", gs)
	}
}

// LevelEnabler 支持通过LevelEnabler接口来处理划分更多详细的日志信息
type LevelEnabler interface {
	Enabled(LogLevel) bool
	Level() LogLevel
}

func (gs LogLevel) Enabled(level LogLevel) bool {
	return level >= gs
}

func (gs LogLevel) Level() LogLevel {
	return gs
}
