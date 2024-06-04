package gslog

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 日志模块
// 1. 日志格式记录(json/text 输出格式)
// 2. 日志元数据(日志级别、时间、文件)
// 3. 日志输出路径(console/file/network)

var (
	_calledPath = 2
)

type Logger struct {
	enabler      LevelEnabler
	writeSyncers map[string]WriteSyncer

	mutex  sync.Mutex
	flag   atomic.Int32
	prefix atomic.Pointer[string]
}

func NewLogger(options ...Option) *Logger {
	gs := &Logger{}

	for _, option := range options {
		option.apply(gs)
	}

	return gs
}

func (gs *Logger) Log(level LogLevel, calledPath int, format string, args ...any) {
	if !gs.enabler.Enabled(level) {
		return
	}
	now := time.Now()
	var file string
	var line int
	var ok bool
	var msg string
	if len(format) > 0 {
		if format[len(format)-1] != '\n' {
			msg = fmt.Sprintf(format, args...)
		} else {
			msg = fmt.Sprintf(format, args...)
		}
	}

	_, file, line, ok = runtime.Caller(calledPath)
	if !ok {
		file = "unknown_file"
		line = 0
	}

	entry := &LogEntry{
		Prefix:     gs.Prefix(),
		FormatFlag: gs.Flags(),
		LogLevel:   gs.Level(),
		File:       file,
		Line:       line,
		CreateAt:   now,
		Message:    msg,
	}

	for _, writeSyncer := range gs.writeSyncers {
		writeSyncer.Sync(entry)
	}
}

func (gs *Logger) Trace(format string, args ...any) {
	gs.Log(TraceLevel, _calledPath, format, args...)
}

func (gs *Logger) Debug(format string, args ...any) {
	gs.Log(DebugLevel, _calledPath, format, args...)

}

func (gs *Logger) Info(format string, args ...any) {
	gs.Log(InfoLevel, _calledPath, format, args...)

}

func (gs *Logger) Warn(format string, args ...any) {
	gs.Log(WarnLevel, _calledPath, format, args...)

}
func (gs *Logger) Error(format string, args ...any) {
	gs.Log(ErrorLevel, _calledPath, format, args...)

}
func (gs *Logger) Critical(format string, args ...any) {
	gs.Log(CriticalLevel, _calledPath, format, args...)
}

func (gs *Logger) Prefix() string {
	if ptrPrefix := gs.prefix.Load(); ptrPrefix != nil {
		return *ptrPrefix
	}
	return ""
}

func (gs *Logger) SetPrefix(prefix string) {
	gs.prefix.Store(&prefix)
}

func (gs *Logger) Flags() int {
	return int(gs.flag.Load())
}

func (gs *Logger) SetFlags(flag int) {
	gs.flag.Store(int32(flag))
}

func (gs *Logger) Level() LogLevel {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	return gs.enabler.Level()
}

func (gs *Logger) SetLevelEnabler(enabler LevelEnabler) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.enabler = enabler
}

func (gs *Logger) AppendWriteSyncer(name string, writeSyncer WriteSyncer) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.writeSyncers[name] = writeSyncer
}

func (gs *Logger) RemoveWriteSyncer(name string) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	writeSyncer, ok := gs.writeSyncers[name]
	if !ok {
		return
	}
	delete(gs.writeSyncers, name)
	_ = writeSyncer.Close()
}

func (gs *Logger) AllWriteSyncers() []string {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	res := make([]string, 0, len(gs.writeSyncers))
	for name, _ := range gs.writeSyncers {
		res = append(res, name)
	}

	return res
}
