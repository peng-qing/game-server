package gslog

import (
	"GameServer/types"
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
