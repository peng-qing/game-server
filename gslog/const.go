package gslog

const (
	SerializeArrayBegin      = '['
	SerializeCommaStep       = ','
	SerializeArrayEnd        = ']'
	SerializePrefixBegin     = '<'
	SerializePrefixEnd       = '>'
	SerializeRadixPointSplit = '.'
	SerializeSpaceSplit      = ' '
	SerializeColonSplit      = ':'
	SerializeFieldStep       = '='
	SerializeNewLine         = '\n'
	SerializeJsonStart       = '{'
	SerializeJsonEnd         = '}'
	SerializeStringMarks     = '"'
)

const (
	badFieldsKey = "!badFieldsKey"
	unknownFile  = "!unknownFile"
)

const (
	BitTextDate            = 1 << iota // 日期标记位
	BitTextTime                        // 时间标记位
	BitTextMicroSecond                 // 微妙标记位
	BitTextFile                        // 文件路径标记位
	BitTextFunction                    // 调用函数标记位
	BitTextLogLevel                    // 日志级别 首字母大写 Trace/Debug/...
	BitTextLogLevelUpCase              // 日志级别 全大写 TRACE/DEBUG/...
	BitTextLogLevelLowCase             // 日志级别 全小写 debug/info/...

	DefaultBitFlag = BitTextDate | BitTextTime | BitTextMicroSecond | BitTextFile | BitTextLogLevel
)

const (
	JsonTimeKey    = "time"
	JsonSourceKey  = "source"
	JsonLevelKey   = "level"
	JsonMessageKey = "message"
	JsonFieldsKey  = "fields"
)
