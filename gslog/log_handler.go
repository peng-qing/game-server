package gslog

import (
	"context"
	"fmt"
	"io"
	"sync"
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
	var err error
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// _ 强迫症
	// 前缀
	if gs.opts.textPrefix != "" {
		_ = buffer.WriteByte(SerializePrefixBegin)
		_, _ = buffer.WriteString(gs.opts.textPrefix)
		_ = buffer.WriteByte(SerializePrefixEnd)
		_ = buffer.WriteByte(SerializeStepSplit)
		// <prefix>
	}

	// 日期
	if gs.opts.flag&BitDate != 0 {
		year, month, day := entry.Time.Date()
		_, _ = buffer.WriteString(fmt.Sprintf("%04d/%02d/%02d", year, month, day))
		_ = buffer.WriteByte(SerializeStepSplit)
		// <prefix> 2024/06/11
	}

	// 时间
	if gs.opts.flag&BitTime != 0 {
		hour, minute, second := entry.Time.Clock()
		_, _ = buffer.WriteString(fmt.Sprintf("%02d:%02d:%02d", hour, minute, second))
		// <prefix> 2024/06/11 10:00:00
		if gs.opts.flag&BitMicroSecond != 0 {
			_ = buffer.WriteByte(SerializeTimeMicroSecondSplit)
			microSec := entry.Time.Nanosecond() / 1000
			_, _ = buffer.WriteString(fmt.Sprintf("%06d", microSec))
		}
		_ = buffer.WriteByte(SerializeStepSplit)
		// <prefix> 2024/06/11 10:00:00.000000
	}

	if gs.opts.flag&(BitLogLevel|BitLogLevelUpCase|BitLogLevelLowCase) != 0 {
		logLevel := gs.opts.level.Level()
		_ = buffer.WriteByte(SerializeArrayBegin)
		if gs.opts.flag&BitLogLevel != 0 {
			_, _ = buffer.WriteString(logLevel.CapitalString())
		} else if gs.opts.flag&BitLogLevelUpCase != 0 {
			_, _ = buffer.WriteString(logLevel.UpCaseString())
		} else {
			_, _ = buffer.WriteString(logLevel.LowCaseString())
		}
		_ = buffer.WriteByte(SerializeArrayEnd)
		_ = buffer.WriteByte(SerializeStepSplit)
		// <prefix> 2024/06/11 10:00:00.000000 [Info]
	}

	if gs.opts.flag&(BitFile|BitFunction) != 0 {
		file, line, function := entry.Source()
		if gs.opts.flag&(BitFile) != 0 {
			if file == "" {
				file = unknownFile
			}
			_, _ = buffer.WriteString(file)
			_ = buffer.WriteByte(SerializeFileLineSplit)
			buffer.AppendInt(int64(line))
			_ = buffer.WriteByte(SerializeStepSplit)
			//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line
		}
		if gs.opts.flag&BitFunction != 0 {
			_, _ = buffer.WriteString(function)
			_ = buffer.WriteByte(SerializeStepSplit)
			//  <prefix> 2024/06/11 10:00:00.000000 [Info] file:line function
		}
	}

	// message
	_, _ = buffer.WriteString(entry.Msg)
	// fields

	_, err = gs.writeSyncer.Write(buffer.Bytes())
	return err
}
