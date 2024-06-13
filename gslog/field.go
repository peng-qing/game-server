package gslog

import (
	"GameServer/types"
	"encoding"
	"encoding/json"
	"reflect"
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

////////// implements

// MarshalText 序列化文本格式 key=value
// 实现 encoding.TextMarshaler 接口
func (gs Field) MarshalText() ([]byte, error) {
	buffer := Get()
	defer buffer.Free()

	buffer.AppendString(gs.Key)
	if gs.Value.Kind() == FieldValueKindField {
		buffer.AppendByte(SerializeRadixPointSplit)
		data, err := gs.Value.Field().MarshalText()
		if err != nil {
			return nil, err
		}
		buffer.AppendBytes(data)
		return buffer.Bytes(), nil
	}
	buffer.AppendByte(SerializeFieldStep)
	if gs.Value.Kind() == FieldValueKindAny {
		if vv, ok := gs.Value.Any().(encoding.TextMarshaler); ok {
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

// MarshalJSON 序列化Json格式 {"key":value}
// 实现 json.Marshaler 接口
func (gs Field) MarshalJSON() ([]byte, error) {
	buffer := Get()
	defer buffer.Free()
	// key
	buffer.AppendByte(SerializeJsonStart)
	buffer.AppendByte(SerializeStringMarks)
	buffer.AppendString(gs.Key)
	buffer.AppendByte(SerializeStringMarks)
	buffer.AppendByte(SerializeColonSplit)

	switch gs.Value.Kind() {
	case FieldValueKindField:
		data, err := gs.Value.Field().MarshalJSON()
		if err != nil {
			return nil, err
		}
		buffer.AppendBytes(data)
	case FieldValueKindAny:
		if vv, ok := gs.Value.Any().(json.Marshaler); ok {
			data, err := vv.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buffer.AppendBytes(data)
			break
		}
		if vv := reflect.ValueOf(gs.Value.Any()); vv.Kind() == reflect.Pointer {
			if vv.IsNil() {
				buffer.AppendByte(SerializeStringMarks)
				buffer.AppendByte(SerializeStringMarks)
				break
			}
			// 先解引用
			elem := vv.Elem()
			data, err := json.Marshal(elem.Interface())
			if err != nil {
				return nil, err
			}
			buffer.AppendBytes(data)
		}
	default:
		buffer.AppendByte(SerializeStringMarks)
		buffer.AppendString(gs.Value.String())
		buffer.AppendByte(SerializeStringMarks)
	}

	buffer.AppendByte(SerializeJsonEnd)
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
