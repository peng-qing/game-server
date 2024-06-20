## 日志


#### 期望

1. 支持结构化日志
2. 对于不同的WriteSyncer支持不同的LevelEnabler
3. 优化日志格式化处理方式


#### 实现

日志支持结构化日志。所有的日志信息都通过 LogHandler 处理。

其定义了对于对应 `LevelEnabler` 的支持和对于日志元数据 `LogEntry` 的处理接口。

```go
type LogHandler interface {
	// Enabled 针对每个Handler支持处理不同的日志级别
	Enabled(ctx context.Context, levelEnabler LevelEnabler) bool
	// LogRecord 处理每条日志元数据
	LogRecord(ctx context.Context, entry *LogEntry) error
}
```

如果需要自定义其他的日志处理方式(比如输出到网络接口)，只需要实现该接口。

默认存在两个 `LogHandler` 的实现方式，分别以 文本格式(`TextHandler`) 和 JSON格式(`JsonHandler`)。

其内部都继承一个内部基础实现`commonHandler`并且存在一个 `io.Writer` 用于最终处理解构化后的日志结果。

```go
type commonHandler struct {
	mutex       sync.Mutex
	writeSyncer io.Writer
	opts        *LogHandlerOptions
}
```

(理论上，Options应该做的更通用)

对于日志参数，可选部分都会最终构成一组 `Fields`。可选参数会依次按照输入数量最终构成不同的 `Field`（最好保证key是string，否则可能出现 ）。

```go
type Field struct {
	Key   string
	Value FieldValue
}
```

`FieldValue` 是结构化的核心。本质上是包含类型值和类型说明，避免过于频繁的使用反射导致的性能下降。

```
type FieldValue struct {
	// 禁止对FieldValue使用 ==
	_ [0]func()
	// FieldValue实际存储的类型说明
	kind FieldValueKind
	// 实际存储的值
	// int/int8/int16/int32/int64/rune ==> int64
	// uint/uint8/uint32/uint64/byte ==> uint64
	// float32/float64 ==> float64
	value any
}
```
