package gslog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"GameServer/common/pool"
)

type LogHandler interface {
	// Enabled 针对每个Handler支持处理不同的日志级别
	Enabled(ctx context.Context, levelEnabler LevelEnabler) bool
	// LogRecord 处理每条日志元数据
	LogRecord(ctx context.Context, entry *LogEntry) error
}

type commonHandler struct {
	mutex       sync.Mutex
	writeSyncer io.Writer
	opts        *LogHandlerOptions
}

func newCommonHandler(writeSyncer io.Writer, opts ...Options) *commonHandler {
	options := &LogHandlerOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}

	return &commonHandler{
		writeSyncer: writeSyncer,
		opts:        options,
	}
}

func (gs *commonHandler) Enabled(_ context.Context, levelEnabler LevelEnabler) bool {
	if gs.opts == nil {
		return TraceLevel.Enabled(levelEnabler.Level())
	}

	return gs.opts.level.Enabled(levelEnabler.Level())
}

func (gs *commonHandler) LogRecord(_ context.Context, entry *LogEntry) error {
	return nil
}

type TextHandler struct {
	*commonHandler
}

func NewTextHandler(writeSyncer io.Writer, opts ...Options) *TextHandler {
	return &TextHandler{
		newCommonHandler(writeSyncer, opts...),
	}
}

func (gs *TextHandler) LogRecord(_ context.Context, entry *LogEntry) error {
	buffer := Get()
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// 前缀
	if gs.opts.textPrefix != "" {
		buffer.AppendByte(SerializePrefixBegin)
		buffer.AppendString(gs.opts.textPrefix)
		buffer.AppendByte(SerializePrefixEnd)
		buffer.AppendByte(SerializeSpaceSplit)
		// <prefix><space>
	}

	// 日期
	if gs.opts.textFlag&BitTextDate != 0 {
		year, month, day := entry.Time.Date()
		buffer.AppendString(fmt.Sprintf("%04d/%02d/%02d", year, month, day))
		buffer.AppendByte(SerializeSpaceSplit)
		// <prefix> 2024/06/11<space>
	}

	// 时间
	if gs.opts.textFlag&BitTextTime != 0 {
		hour, minute, second := entry.Time.Clock()
		buffer.AppendString(fmt.Sprintf("%02d:%02d:%02d", hour, minute, second))
		// <prefix> 2024/06/11 10:00:00
		if gs.opts.textFlag&BitTextMicroSecond != 0 {
			buffer.AppendByte(SerializeRadixPointSplit)
			microSec := entry.Time.Nanosecond() / 1000
			buffer.AppendString(fmt.Sprintf("%06d", microSec))
		}
		buffer.AppendByte(SerializeSpaceSplit)
		// <prefix> 2024/06/11 10:00:00.000000<space>
	}

	if gs.opts.textFlag&(BitTextLogLevel|BitTextLogLevelUpCase|BitTextLogLevelLowCase) != 0 {
		logLevel := entry.Level
		buffer.AppendByte(SerializeArrayBegin)
		if gs.opts.textFlag&BitTextLogLevel != 0 {
			buffer.AppendString(logLevel.CapitalString())
		} else if gs.opts.textFlag&BitTextLogLevelUpCase != 0 {
			buffer.AppendString(logLevel.UpCaseString())
		} else {
			buffer.AppendString(logLevel.LowCaseString())
		}
		buffer.AppendByte(SerializeArrayEnd)
		buffer.AppendByte(SerializeSpaceSplit)
		// <prefix> 2024/06/11 10:00:00.000000 [Info]<space>
	}

	if gs.opts.textFlag&(BitTextFile|BitTextFunction) != 0 {
		file, line, function := entry.Source()
		if gs.opts.textFlag&(BitTextFile) != 0 {
			if file == "" {
				file = unknownFile
			}
			buffer.AppendString(file)
			buffer.AppendByte(SerializeColonSplit)
			buffer.AppendInt(int64(line))
			buffer.AppendByte(SerializeSpaceSplit)
			//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line<space>
		}
		if gs.opts.textFlag&BitTextFunction != 0 {
			buffer.AppendString(function)
			buffer.AppendByte(SerializeSpaceSplit)
			//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line function<space>
		}
	}

	// message
	{
		buffer.AppendString(entry.Msg)
		buffer.AppendByte(SerializeSpaceSplit)
		//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line function<space>message<space>
	}
	// fields
	for _, field := range entry.Fields {
		data, err := field.SerializeText()
		if err != nil {
			continue
		}
		buffer.AppendBytes(data)
		buffer.AppendByte(SerializeSpaceSplit)
		//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line function<space>message fieldKey=fieldValue<space>...
	}
	buffer.AppendByte(SerializeNewLine)
	//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line function<space>message fieldKey=fieldValue<space>...\n

	_, err := gs.writeSyncer.Write(buffer.Bytes())
	return err
}

type JsonHandler struct {
	*commonHandler
}

func NewJsonHandler(writeSyncer io.Writer, opts ...Options) *JsonHandler {
	return &JsonHandler{
		newCommonHandler(writeSyncer, opts...),
	}
}

func (gs *JsonHandler) LogRecord(_ context.Context, entry *LogEntry) error {
	buffer := Get()
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	buffer.AppendByte(SerializeJsonStart)
	// 时间
	{
		if !entry.Time.IsZero() {
			format := gs.opts.jsonTimeFormat
			if format == "" {
				format = time.RFC3339
			}
			gs.appendJsonKey(buffer, JsonTimeKey)
			gs.appendJsonValue(buffer, entry.Time.Format(format))
		}
	}
	// 来源
	{
		gs.appendJsonKey(buffer, JsonSourceKey)
		file, line, function := entry.Source()
		sourceStr := fmt.Sprintf("%s:%d %s", file, line, function)
		gs.appendJsonValue(buffer, sourceStr)
	}
	// 日志级别
	{
		gs.appendJsonKey(buffer, JsonLevelKey)
		gs.appendJsonValue(buffer, entry.Level.LowCaseString())
	}
	// msg
	{
		gs.appendJsonKey(buffer, JsonMessageKey)
		gs.appendJsonValue(buffer, entry.Msg)
	}
	// fields...
	{
		for _, field := range entry.Fields {
			gs.appendJsonKey(buffer, field.Key)
			gs.appendJsonValue(buffer, field.Value)

			//gs.appendJsonValue(buffer, field.Value.String())
			//	data, err := field.SerializeJson()
			//	if err != nil {
			//		continue
			//	}
			//	buffer.AppendBytes(data)
		}
	}
	buffer.AppendByte(SerializeJsonEnd)
	buffer.AppendByte(SerializeNewLine)

	_, err := gs.writeSyncer.Write(buffer.Bytes())

	return err
}

func (gs *JsonHandler) appendJsonKey(buffer *pool.Buffer, key string) {
	/// json string,
	if buffer.Size() != 1 {
		buffer.AppendByte(SerializeCommaStep)
	}
	// "key":<space>
	buffer.AppendByte(SerializeStringMarks)
	buffer.AppendString(key)
	buffer.AppendByte(SerializeStringMarks)
	buffer.AppendByte(SerializeColonSplit)
	buffer.AppendByte(SerializeSpaceSplit)
}

// TODO 需要重构 目前对于各种类型支持的不好 并且如果是Field无法处理
func (gs *JsonHandler) appendJsonValue(buffer *pool.Buffer, val any) {
	defer func() {
		if err := recover(); err != nil {
			buffer.AppendByte(SerializeStringMarks)
			_ = fmt.Sprintf("%+v", val)
			buffer.AppendByte(SerializeStringMarks)
		}
	}()

	// 实现 json.Marshaler 接口
	if vv, ok := val.(json.Marshaler); ok {
		data, err := vv.MarshalJSON()
		if err != nil {
			panic(err)
		}
		buffer.AppendBytes(data)
		return
	}
	// slice
	//if reflect.TypeOf(val).Kind() == reflect.Slice {
	//	sliceVals := reflect.ValueOf(val)
	//	res := make([]any, sliceVals.Len())
	//	for i := 0; i < sliceVals.Len(); i++ {
	//		res[i] = sliceVals.Index(i).Interface()
	//	}
	//	encBuf := Get()
	//	enc := json.NewEncoder(encBuf)
	//	enc.SetEscapeHTML(false)
	//	if err := enc.Encode(res); err != nil {
	//		panic(err)
	//	}
	//	encBuf.TrimNewLine()
	//	buffer.AppendBytes(encBuf.Bytes())
	//	return
	//}

	switch vv := val.(type) {
	case string:
		buffer.AppendByte(SerializeStringMarks)
		buffer.AppendString(vv)
		buffer.AppendByte(SerializeStringMarks)
	case byte:
		buffer.AppendByte(SerializeStringMarks)
		buffer.AppendByte(vv)
		buffer.AppendByte(SerializeStringMarks)
	case time.Time:
		buffer.AppendByte(SerializeStringMarks)
		if !vv.IsZero() {
			format := gs.opts.jsonTimeFormat
			if format == "" {
				format = time.RFC3339
			}
			buffer.AppendTime(vv, format)
		}
		buffer.AppendByte(SerializeStringMarks)
	case bool:
		buffer.AppendBool(vv)
	case int, int8, int32, int16, int64:
		buffer.AppendInt(vv.(int64))
	case uint, uint16, uint32, uint64:
		buffer.AppendUint(vv.(uint64))
	case float32, float64:
		buffer.AppendFloat(vv.(float64), 64)
	default:
		buffer.AppendByte(SerializeStringMarks)
		_ = fmt.Sprintf("%+v", val)
		buffer.AppendByte(SerializeStringMarks)
	}
}
