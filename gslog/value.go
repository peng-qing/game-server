package gslog

import (
	"GameServer/utils"
	"fmt"
	"strconv"
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
)

var kindString = map[FieldValueKind]string{
	FieldValueKindAny:      "Any",
	FieldValueKindInt64:    "Int64",
	FieldValueKindInt64s:   "Int64Array",
	FieldValueKindUint64:   "Uint64",
	FieldValueKindUint64s:  "Uint64Array",
	FieldValueKindFloat64:  "Float64",
	FieldValueKindFloat64s: "Float64Array",
	FieldValueKindString:   "String",
	FieldValueKindBool:     "Bool",
	FieldValueKindBools:    "BoolArray",
	FieldValueKindTime:     "Time",
	FieldValueKindDuration: "Duration",
	FieldValueKindField:    "Field",
	FieldValueKindFields:   "FieldArray",
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

	// TODO 优化根据类型处理
	buffer := make([]byte, 0)
	fmt.Append(buffer, gs.value)

	return string(buffer)
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

////////// internal

// format to dst... like fmt.Sprint
func (gs FieldValue) append(dst []byte) []byte {
	switch gs.Kind() {
	case FieldValueKindAny:
		fmt.Append(dst, gs.value)
	case FieldValueKindInt64:
		strconv.AppendInt(dst, gs.value.(int64), 10)
	case FieldValueKindUint64:
		strconv.AppendUint(dst, gs.value.(uint64), 10)
	case FieldValueKindFloat64:
		strconv.AppendFloat(dst, gs.value.(float64), 'f', -1, 64)
	case FieldValueKindBool:
		strconv.AppendBool(dst, gs.value.(bool))
	case FieldValueKindTime:
		return append(dst, gs.value.(time.Time).String()...)
	case FieldValueKindDuration:
		return append(dst, gs.value.(time.Duration).String()...)
	case FieldValueKindField:
	case FieldValueKindFields:
	case FieldValueKindInt64s:
	case FieldValueKindUint64s:
	case FieldValueKindFloat64s:
	case FieldValueKindStrings:
	case FieldValueKindBools:
	default:
		panic(fmt.Sprintf("Invalid FieldValueKind %s", kindString[gs.Kind()]))
	}
	return dst
}
