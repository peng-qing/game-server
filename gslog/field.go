package gslog

import (
	"GameServer/types"
	"encoding"
	"time"
)

type Field struct {
	Key   string
	Value FieldValue
}

////// constructs

func Sting[T types.String](key string, val T) Field {
	return Field{
		Key:   key,
		Value: StringFieldValue(string(val)),
	}
}

func Int[T types.Int](key string, val T) Field {
	return Field{
		Key:   key,
		Value: Int64FieldValue(int64(val)),
	}
}

func Uint[T types.Uint](key string, val T) Field {
	return Field{
		Key:   key,
		Value: Uint64FieldValue(uint64(val)),
	}
}

func Bool[T types.Bool](key string, val T) Field {
	return Field{
		Key:   key,
		Value: BoolFieldValue(bool(val)),
	}
}

func Float[T types.Float](key string, val T) Field {
	return Field{
		Key:   key,
		Value: Float64FieldValue(float64(val)),
	}
}

func Errors(key string, val ...error) Field {
	return Field{
		Key:   key,
		Value: ErrorFieldValue(val...),
	}
}

func Fields(key string, val ...Field) Field {
	fields := Field{Key: key}
	if len(val) <= 1 {
		fields.Value = FieldFieldValue(val[0])
		return fields
	}
	fields.Value = FieldArrayFieldValue(val...)
	return fields
}

func Time(key string, val time.Time) Field {
	return Field{
		Key:   key,
		Value: TimeFieldValue(val),
	}
}

func Duration(key string, val time.Duration) Field {
	return Field{
		Key:   key,
		Value: DurationFieldValue(val),
	}
}

func Any(key string, val any) Field {
	return Field{
		Key:   key,
		Value: AnyFieldValue(val),
	}
}

////////// Accessors

// SerializeText 序列化文本格式 key=value
func (gs Field) SerializeText() ([]byte, error) {
	buffer := Get()
	_, _ = buffer.WriteString(gs.Key)
	if gs.Value.Kind() == FieldValueKindField {
		buffer.AppendByte(SerializeRadixPointSplit)
		data, err := gs.Value.Field().SerializeText()
		if err != nil {
			return nil, err
		}
		buffer.AppendBytes(data)
		return buffer.Bytes(), nil
	}
	buffer.AppendByte(SerializeFieldStep)
	if gs.Value.Kind() == FieldValueKindAny {
		if vv, ok := gs.Value.value.(encoding.TextMarshaler); ok {
			data, err := vv.MarshalText()
			if err != nil {
				return nil, err
			}
			buffer.AppendBytes(data)
			return buffer.Bytes(), nil
		}
	}
	// default
	buffer.AppendString(gs.Value.String())

	return buffer.Bytes(), nil
}

// SerializeJson 序列化json格式 key: value
func (gs Field) SerializeJson() ([]byte, error) {
	buffer := Get()
	buffer.AppendByte(SerializeRadixPointSplit)

	return buffer.Bytes(), nil
}

////////// internal

func argsToFields(args ...any) (Field, []any) {
	switch v := args[0].(type) {
	case Field:
		return v, args[1:]
	case string:
		if len(args) == 1 {
			return Sting[string](badFieldsKey, v), nil
		}
		return Any(v, args[1]), args[2:]
	default:
		return Any(badFieldsKey, v), args[1:]
	}
}
