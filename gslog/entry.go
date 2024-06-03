package gslog

import (
	"time"
)

//	日志信息标记位
//
// bitmap 形式, 允许自定义相关输出形式
const (
	BitDate            = 1 << iota // 日期标记位 2024/06/01
	BitTime                        // 时间标记位 10:00:00
	BitMicroSeconds                // 微秒标记位 10:00:00.000000
	BitFullPath                    // 完整文件名标记位	/home/go/project/main.go
	BitShortFile                   // 短文件名标记位 main.go
	BitLogLevel                    // 日志级别标记位 Debug/Info/Warn/Error....
	BitLogLevelUpCase              // 日志级别 大写 DEBUG/INFO/....
	BitLogLevelLowCase             // 日志级别 小写 debug/info/warn/....
)

const (
	BitDefaultStdSourceFlag = BitDate | BitTime | BitMicroSeconds | BitFullPath | BitLogLevel
)

type LogEntry struct {
	// 来自Logger的信息
	Prefix     string
	FormatFlag int
	// 日志元数据
	File     string
	Line     int
	CreateAt time.Time
	Message  string
}

//func (gs *LogEntry) Format() []byte {
//
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
//}
