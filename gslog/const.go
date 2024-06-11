package gslog

const (
	SerializeArrayBegin = '['
	SerializeArrayStep  = ','
	SerializeArrayEnd   = ']'

	SerializePrefixBegin = '<'
	SerializePrefixEnd   = '>'

	SerializeTimeMicroSecondSplit = '.'

	SerializeStepSplit = ' '

	SerializeFileLineSplit = ':'
)

const (
	badFieldsKey = "!badFieldsKey"
	unknownFile  = "!unknownFile"
)

const (
	BitDate            = 1 << iota // 日期标记位
	BitTime                        // 时间标记位
	BitMicroSecond                 // 微妙标记位
	BitFile                        // 文件路径标记位
	BitFunction                    // 调用函数标记位
	BitLogLevel                    // 日志级别 首字母大写 Trace/Debug/...
	BitLogLevelUpCase              // 日志级别 全大写 TRACE/DEBUG/...
	BitLogLevelLowCase             // 日志级别 全小写 debug/info/...
)
