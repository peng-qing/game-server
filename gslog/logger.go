package gslog

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 日志模块
// 1. 日志格式记录(json/text 输出格式)
// 2. 日志元数据(日志级别、时间、文件)
// 3. 日志输出路径(console/file/network)

type Logger struct {
	enabler      LevelEnabler
	writeSyncers []WriteSyncer

	mu sync.Mutex

	formatFlag int
	calledPath int
	prefix     string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (gs *Logger) Log(level LogLevel, format string, args ...any) {
	if !gs.enabler.Enabled(level) {
		return
	}
	now := time.Now()
	var file string
	var line int
	var ok bool
	//buffer := Get()
	var msg string
	if len(format) > 0 {
		if format[len(format)-1] != '\n' {
			msg = fmt.Sprintln(format, args)
		} else {
			msg = fmt.Sprintf(format, args)
		}
	}

	_, file, line, ok = runtime.Caller(gs.calledPath)
	if !ok {
		file = "unknown_file"
		line = 0
	}

	entry := &LogEntry{
		Prefix:     gs.prefix,
		FormatFlag: gs.formatFlag,
		File:       file,
		Line:       line,
		CreateAt:   now,
		Message:    msg,
	}

	for _, writeSyncer := range gs.writeSyncers {
		writeSyncer.Sync(entry)
	}

	//gs.mu.Lock()
	//defer gs.mu.Unlock()
	// 不在这格式化 丢给写的时候格式化
	//if gs.formatFlag&(BitFullPath|BitShortFile) != 0 {
	//	// 需要输出文件信息
	//	var ok bool
	//	_, file, line, ok = runtime.Caller(gs.calledPath)
	//	if !ok {
	//		file = "unknown_file"
	//		line = 0
	//	}
	//}
	//// 1. header string
	//if gs.prefix != "" {
	//	// 强迫症....
	//	_, _ = buffer.WriteString(fmt.Sprintf("<%s> ", gs.prefix))
	//	// <prefix>
	//}
	//if gs.formatFlag&BitDate != 0 {
	//	year, month, day := now.Date()
	//	_, _ = buffer.WriteString(fmt.Sprintf("%04d/%02d/%02d ", year, month, day))
	//	// <prefix> 2024/06/01
	//}
	//if gs.formatFlag&BitTime != 0 {
	//	hour, minute, second := now.Clock()
	//	_, _ = buffer.WriteString(fmt.Sprintf("%02d:%02d:%02d", hour, minute, second))
	//	// <prefix> 2024/06/01 00:00:00
	//}
	//if gs.formatFlag&BitMicroSeconds != 0 {
	//	microSec := now.Nanosecond() / 1000
	//	_, _ = buffer.WriteString(fmt.Sprintf(".%06d ", microSec))
	//	// <prefix> 2024/06/01 00:00:00.000000
	//}
	//if gs.formatFlag&(BitLogLevel|BitLogLevelUpCase|BitLogLevelLowCase) != 0 {
	//	logLevel := gs.enabler.Level()
	//	if gs.formatFlag&(BitLogLevelUpCase|BitLogLevelLowCase) == 0 {
	//		_, _ = buffer.WriteString(fmt.Sprintf("[%s] ", logLevel.CapitalString()))
	//	} else if gs.formatFlag&BitLogLevelUpCase != 0 {
	//		_, _ = buffer.WriteString(fmt.Sprintf("[%s] ", logLevel.UpCaseString()))
	//	} else {
	//		_, _ = buffer.WriteString(fmt.Sprintf("[%s] ", logLevel.LowCaseString()))
	//	}
	//	// <prefix> 2024/06/01 00:00:00.000000 [Info]
	//}
	//if gs.formatFlag&(BitShortFile|BitFullPath) != 0 {
	//	if gs.formatFlag&BitShortFile != 0 {
	//		for i := len(file); i >= 0; i-- {
	//			if file[i] == '/' {
	//				file = file[i+1:]
	//				break
	//			}
	//		}
	//	}
	//	_, _ = buffer.WriteString(fmt.Sprintf("%s:%d  ", file, line))
	//	// <prefix> 2024/06/01 00:00:00.000000 [Info] file:line
	//}

	//gs.encoder.EncoderEntry(level, format, args)
	////for _, syncer := range gs.writeSyncers {
	////	//syncer.Write()
	////}
}

func (gs *Logger) Trace(format string, args ...any) {
	gs.Log(TraceLevel, format, args)
}

func (gs *Logger) Debug(format string, args ...any) {
	gs.Log(DebugLevel, format, args)

}

func (gs *Logger) Info(format string, args ...any) {
	gs.Log(InfoLevel, format, args)

}

func (gs *Logger) Warn(format string, args ...any) {
	gs.Log(WarnLevel, format, args)

}
func (gs *Logger) Error(format string, args ...any) {
	gs.Log(ErrorLevel, format, args)

}
func (gs *Logger) Critical(format string, args ...any) {
	gs.Log(CriticalLevel, format, args)
}
