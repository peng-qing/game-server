package v1

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	errUnmarshalNilPointer = errors.New("unmarshal text to nil LogLevel pointer")
	errUnmarshalInvalid    = errors.New("unmarshal invalid text")
)

var (
	_ LevelEnabler = (*LogLevel)(nil)
)

// LevelEnabler 用于判断日志级别是否开启
type LevelEnabler interface {
	Enabled(level LogLevel) bool
	Level() LogLevel
}

type LogLevel int8

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel

	_minLevel = TraceLevel
	_maxLevel = CriticalLevel

	// InvalidLevel 无效值
	InvalidLevel = TraceLevel - 1
)

// ParseLogLevel 解析字符串到LogLevel
func ParseLogLevel(text string) (LogLevel, error) {
	var lv LogLevel
	err := lv.UnmarshalText([]byte(text))

	return lv, err
}

// UnmarshalText 解析文本为指定日志级别
func (gs *LogLevel) UnmarshalText(text []byte) error {
	if gs == nil {
		// 指针接收器 有可能空
		return errUnmarshalNilPointer
	}
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

func (gs LogLevel) Enabled(level LogLevel) bool {
	if level < _minLevel || level > _maxLevel {
		return false
	}
	return level >= gs
}

func (gs LogLevel) Level() LogLevel {
	return gs
}
