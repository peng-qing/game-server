package gslog

import (
	"runtime"
	"time"
)

type LogEntry struct {
	Time   time.Time
	Level  LogLevel
	PC     uintptr
	Msg    string
	Fields []Field
}

func NewLogEntry(t time.Time, level LogLevel, msg string, pc uintptr) *LogEntry {
	return &LogEntry{
		Time:  t,
		Level: level,
		PC:    pc,
		Msg:   msg,
	}
}

func (gs *LogEntry) AppendFields(fields ...Field) {
	gs.Fields = append(gs.Fields, fields...)
}

func (gs *LogEntry) Source() (file string, line int, function string) {
	frames := runtime.CallersFrames([]uintptr{gs.PC})
	frame, _ := frames.Next()

	return frame.File, frame.Line, frame.Function
}

func (gs *LogEntry) AddArgs(args ...any) {
	var field Field
	for len(args) > 0 {
		field, args = argsToFields(args...)
		gs.AppendFields(field)
	}
}
