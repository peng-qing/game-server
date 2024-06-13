package gslog

import (
	"GameServer/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 如果要实现结构化日志
// 1. 如何记录结构化数据
// 2. 以何种形式格式化结构化数据

// 为了方便展示 Value需要对特殊类型进行处理
// 1.基础数据类型
// 2.Field类型，表示出现层级，需要层级处理展示

type FieldValueKind int

const (
	FieldValueKindAny FieldValueKind = iota
	FieldValueKindInt64
	FieldValueKindInt64s
	FieldValueKindUint64
	FieldValueKindUint64s
	FieldValueKindFloat64
	FieldValueKindFloat64s
	FieldValueKindString
	FieldValueKindStrings
	FieldValueKindBool
	FieldValueKindBools
	FieldValueKindTime
	FieldValueKindDuration
	FieldValueKindField
	FieldValueKindFields
	FieldValueKindError
)

var kindString = map[FieldValueKind]string{
	FieldValueKindAny:      "Any",
	FieldValueKindInt64:    "Int64",
	FieldValueKindInt64s:   "Int64s",
	FieldValueKindUint64:   "Uint64",
	FieldValueKindUint64s:  "Uint64s",
	FieldValueKindFloat64:  "Float64",
	FieldValueKindFloat64s: "Float64s",
	FieldValueKindString:   "String",
	FieldValueKindBool:     "Bool",
	FieldValueKindBools:    "Bools",
	FieldValueKindTime:     "Time",
	FieldValueKindDuration: "Duration",
	FieldValueKindField:    "Field",
	FieldValueKindFields:   "Fields",
	FieldValueKindError:    "Error",
}

func (gs FieldValueKind) String() string {
	if gs >= 0 && int(gs) < len(kindString) {
		return kindString[gs]
	}
	return "UnknownFieldValueKind"
}

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

/////////////////  FieldValue constructs

func IntFieldValue(val int) FieldValue {
	return FieldValue{kind: FieldValueKindInt64, value: int64(val)}
}

func Int64FieldValue(val int64) FieldValue {
	return FieldValue{kind: FieldValueKindInt64, value: val}
}

func Int64ArrayFieldValue(val ...int64) FieldValue {
	return FieldValue{kind: FieldValueKindInt64s, value: val}
}

func Uint64FieldValue(val uint64) FieldValue {
	return FieldValue{kind: FieldValueKindUint64, value: val}
}

func Uint64ArrayFieldValue(val ...uint64) FieldValue {
	return FieldValue{kind: FieldValueKindUint64s, value: val}
}

func Float64FieldValue(val float64) FieldValue {
	return FieldValue{kind: FieldValueKindFloat64, value: val}
}

func Float64ArrayFieldValue(val ...float64) FieldValue {
	return FieldValue{kind: FieldValueKindFloat64s, value: val}
}

func StringFieldValue(val string) FieldValue {
	return FieldValue{kind: FieldValueKindString, value: val}
}

func StringArrayFieldValue(val ...string) FieldValue {
	return FieldValue{kind: FieldValueKindStrings, value: val}
}

func BoolFieldValue(val bool) FieldValue {
	return FieldValue{kind: FieldValueKindBool, value: val}
}

func BoolArrayFieldValue(val ...bool) FieldValue {
	return FieldValue{kind: FieldValueKindBools, value: val}
}

func TimeFieldValue(val time.Time) FieldValue {
	return FieldValue{kind: FieldValueKindTime, value: val}
}

func DurationFieldValue(val time.Duration) FieldValue {
	return FieldValue{kind: FieldValueKindDuration, value: val}
}

func FieldFieldValue(val Field) FieldValue {
	return FieldValue{kind: FieldValueKindField, value: val}
}

func FieldArrayFieldValue(val ...Field) FieldValue {
	return FieldValue{kind: FieldValueKindFields, value: val}
}

func ErrorFieldValue(errs ...error) FieldValue {
	return FieldValue{kind: FieldValueKindError, value: errors.Join(errs...)}
}

func AnyFieldValue(val any) FieldValue {
	switch vv := val.(type) {
	case int:
		return IntFieldValue(vv)
	case []int:
		return Int64ArrayFieldValue(utils.ConvertIntSliceToInt64s[int](vv)...)
	case int8:
		return Int64FieldValue(int64(vv))
	case []int8:
		return Int64ArrayFieldValue(utils.ConvertIntSliceToInt64s[int8](vv)...)
	case int16:
		return Int64FieldValue(int64(vv))
	case []int16:
		return Int64ArrayFieldValue(utils.ConvertIntSliceToInt64s[int16](vv)...)
	case int32:
		return Int64FieldValue(int64(vv))
	case []int32:
		return Int64ArrayFieldValue(utils.ConvertIntSliceToInt64s[int32](vv)...)
	case int64:
		return Int64FieldValue(vv)
	case []int64:
		return Int64ArrayFieldValue(vv...)
	case uint:
		return Uint64FieldValue(uint64(vv))
	case []uint:
		return Uint64ArrayFieldValue(utils.ConvertUintSliceToUint64s[uint](vv)...)
	case uint8:
		return Uint64FieldValue(uint64(vv))
	case []uint8:
		return Uint64ArrayFieldValue(utils.ConvertUintSliceToUint64s[uint8](vv)...)
	case uint16:
		return Uint64FieldValue(uint64(vv))
	case []uint16:
		return Uint64ArrayFieldValue(utils.ConvertUintSliceToUint64s[uint16](vv)...)
	case uint32:
		return Uint64FieldValue(uint64(vv))
	case []uint32:
		return Uint64ArrayFieldValue(utils.ConvertUintSliceToUint64s[uint32](vv)...)
	case uint64:
		return Uint64FieldValue(vv)
	case []uint64:
		return Uint64ArrayFieldValue(vv...)
	case float32:
		return Float64FieldValue(float64(vv))
	case []float32:
		return Float64ArrayFieldValue(utils.ConvertFloatSliceToFloat64s(vv)...)
	case float64:
		return Float64FieldValue(vv)
	case []float64:
		return Float64ArrayFieldValue(vv...)
	case string:
		return StringFieldValue(vv)
	case []string:
		return StringArrayFieldValue(vv...)
	case bool:
		return BoolFieldValue(vv)
	case []bool:
		return BoolArrayFieldValue(vv...)
	case time.Time:
		return TimeFieldValue(vv)
	case time.Duration:
		return DurationFieldValue(vv)
	case Field:
		return FieldFieldValue(vv)
	case []Field:
		return FieldArrayFieldValue(vv...)
	case error:
		return ErrorFieldValue(vv)
	default:
		return FieldValue{kind: FieldValueKindAny, value: val}
	}
}

//////////// Accessors

func (gs FieldValue) Kind() FieldValueKind {
	return gs.kind
}

func (gs FieldValue) Int64() int64 {
	if current, target := gs.Kind(), FieldValueKindInt64; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(int64)
}

func (gs FieldValue) Int64s() []int64 {
	if current, target := gs.Kind(), FieldValueKindInt64s; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.([]int64)
}

func (gs FieldValue) Uint64() uint64 {
	if current, target := gs.Kind(), FieldValueKindUint64; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(uint64)
}

func (gs FieldValue) Uint64s() []uint64 {
	if current, target := gs.Kind(), FieldValueKindUint64s; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.([]uint64)
}

func (gs FieldValue) Float64() float64 {
	if current, target := gs.Kind(), FieldValueKindFloat64; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(float64)
}

func (gs FieldValue) Float64s() []float64 {
	if current, target := gs.Kind(), FieldValueKindFloat64s; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.([]float64)
}

// String never panic
func (gs FieldValue) String() string {
	if current, target := gs.Kind(), FieldValueKindString; current == target {
		return gs.value.(string)
	}

	buffer := make([]byte, 0)
	return string(gs.appendFieldValue(buffer))
}

func (gs FieldValue) Strings() []string {
	if current, target := gs.Kind(), FieldValueKindStrings; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.([]string)
}

func (gs FieldValue) Bool() bool {
	if current, target := gs.Kind(), FieldValueKindBool; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(bool)
}

func (gs FieldValue) Bools() []bool {
	if current, target := gs.Kind(), FieldValueKindBools; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.([]bool)
}

func (gs FieldValue) Time() time.Time {
	if current, target := gs.Kind(), FieldValueKindTime; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(time.Time)
}

func (gs FieldValue) Duration() time.Duration {
	if current, target := gs.Kind(), FieldValueKindDuration; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(time.Duration)
}

func (gs FieldValue) Field() Field {
	if current, target := gs.Kind(), FieldValueKindField; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(Field)
}

func (gs FieldValue) Fields() []Field {
	if current, target := gs.Kind(), FieldValueKindFields; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.([]Field)
}

func (gs FieldValue) Error() error {
	if current, target := gs.Kind(), FieldValueKindError; current != target {
		panic(fmt.Sprintf("current FieldValueKind is %s, not %s", kindString[current], kindString[target]))
	}
	return gs.value.(error)
}

func (gs FieldValue) Any() any {
	switch gs.Kind() {
	case FieldValueKindAny:
		return gs.value
	case FieldValueKindInt64:
		return gs.Int64()
	case FieldValueKindInt64s:
		return gs.Int64s()
	case FieldValueKindUint64:
		return gs.Uint64()
	case FieldValueKindUint64s:
		return gs.Uint64s()
	case FieldValueKindFloat64:
		return gs.Float64()
	case FieldValueKindFloat64s:
		return gs.Float64s()
	case FieldValueKindString:
		return gs.String()
	case FieldValueKindStrings:
		return gs.Strings()
	case FieldValueKindBool:
		return gs.Bool()
	case FieldValueKindBools:
		return gs.Bools()
	case FieldValueKindTime:
		return gs.Time()
	case FieldValueKindDuration:
		return gs.Duration()
	case FieldValueKindField:
		return gs.Field()
	case FieldValueKindFields:
		return gs.Fields()
	case FieldValueKindError:
		return gs.Error()
	default:
		panic(fmt.Sprintf("unknown kind %s", gs.Kind()))
	}
}

////////// internal

// format to dst... like fmt.Sprint
func (gs FieldValue) appendFieldValue(dst []byte) []byte {
	switch gs.Kind() {
	case FieldValueKindInt64:
		return strconv.AppendInt(dst, gs.value.(int64), 10)
	case FieldValueKindUint64:
		return strconv.AppendUint(dst, gs.value.(uint64), 10)
	case FieldValueKindFloat64:
		return strconv.AppendFloat(dst, gs.value.(float64), 'f', -1, 64)
	case FieldValueKindBool:
		return strconv.AppendBool(dst, gs.value.(bool))
	case FieldValueKindTime:
		return gs.value.(time.Time).AppendFormat(dst, time.RFC3339)
	case FieldValueKindDuration:
		return append(dst, gs.value.(time.Duration).String()...)
	case FieldValueKindString:
		return append(dst, gs.value.(string)...)
	case FieldValueKindError:
		return append(dst, fmt.Sprintf("err: %s", gs.value.(error).Error())...)
	case FieldValueKindAny:
		return fmt.Append(dst, gs.value)
	case FieldValueKindField:
		return fmt.Append(dst, gs.value)
	case FieldValueKindFields:
		return append(dst, gs.serializeFields()...)
	case FieldValueKindInt64s:
		return append(dst, gs.serializeInt64s()...)
	case FieldValueKindUint64s:
		return append(dst, gs.serializeUint64s()...)
	case FieldValueKindFloat64s:
		return append(dst, gs.serializeFloat64s()...)
	case FieldValueKindStrings:
		return append(dst, gs.serializeStrings()...)
	case FieldValueKindBools:
		return append(dst, gs.serializeBools()...)
	default:
		panic(fmt.Sprintf("Invalid FieldValueKind %s", kindString[gs.Kind()]))
	}
	return dst
}

func (gs FieldValue) serializeFields() []byte {
	var buffer []byte
	fields, ok := gs.value.([]Field)
	if !ok {
		return buffer
	}
	buffer = append(buffer, SerializeArrayBegin)
	for idx, field := range fields {
		if idx > 0 {
			buffer = append(buffer, SerializeCommaStep, SerializeSpaceSplit)
		}
		data, err := field.MarshalText()
		if err != nil {
			buffer = fmt.Append(buffer, field)
			continue
		}
		buffer = append(buffer, data...)
	}
	return append(buffer, SerializeArrayEnd)
}

func (gs FieldValue) serializeInt64s() []byte {
	buffer := Get()
	defer buffer.Free()

	nums, ok := gs.value.([]int64)
	if !ok {
		return buffer.Bytes()
	}
	buffer.AppendByte(SerializeArrayBegin)
	for idx, num := range nums {
		if idx > 0 {
			buffer.AppendByte(SerializeCommaStep)
			buffer.AppendByte(SerializeSpaceSplit)
		}
		buffer.AppendInt(num)
	}
	buffer.AppendByte(SerializeArrayEnd)

	return buffer.Bytes()
}

func (gs FieldValue) serializeUint64s() []byte {
	buffer := Get()
	defer buffer.Free()

	nums, ok := gs.value.([]uint64)
	if !ok {
		return buffer.Bytes()
	}

	buffer.AppendByte(SerializeArrayBegin)
	for idx, num := range nums {
		if idx > 0 {
			buffer.AppendByte(SerializeCommaStep)
			buffer.AppendByte(SerializeSpaceSplit)
		}
		buffer.AppendUint(num)
	}
	buffer.AppendByte(SerializeArrayEnd)

	return buffer.Bytes()
}

func (gs FieldValue) serializeFloat64s() []byte {
	buffer := Get()
	defer buffer.Free()

	nums, ok := gs.value.([]float64)
	if !ok {
		return buffer.Bytes()
	}

	buffer.AppendByte(SerializeArrayBegin)
	for idx, num := range nums {
		if idx > 0 {
			buffer.AppendByte(SerializeCommaStep)
			buffer.AppendByte(SerializeSpaceSplit)
		}
		buffer.AppendFloat(num, 64)
	}
	buffer.AppendByte(SerializeArrayEnd)

	return buffer.Bytes()
}

// string 可以使用 strings.Builder 构建
func (gs FieldValue) serializeStrings() []byte {
	var builder strings.Builder
	strs, ok := gs.value.([]string)
	if !ok {
		return nil
	}

	builder.WriteByte(SerializeArrayBegin)
	for idx, str := range strs {
		if idx > 0 {
			builder.WriteByte(SerializeCommaStep)
			builder.WriteByte(SerializeSpaceSplit)
		}
		builder.WriteString(str)
	}
	builder.WriteByte(SerializeArrayEnd)

	return []byte(builder.String())
}

func (gs FieldValue) serializeBools() []byte {
	buffer := Get()
	defer buffer.Free()

	bools, ok := gs.value.([]bool)
	if !ok {
		return nil
	}

	buffer.AppendByte(SerializeArrayBegin)
	for idx, boolVal := range bools {
		if idx > 0 {
			buffer.AppendByte(SerializeCommaStep)
			buffer.AppendByte(SerializeSpaceSplit)
		}
		buffer.AppendBool(boolVal)
	}
	buffer.AppendByte(SerializeArrayEnd)

	return buffer.Bytes()
}
