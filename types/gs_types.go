package types

import "cmp"

type None struct{}

type Integer interface {
	Int
	Uint
	Character
	Pointer
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type String interface {
	~string | ~[]byte | ~[]rune
}

type Character interface {
	~byte | ~rune
}

type Float interface {
	~float32 | ~float64
}

type Complex interface {
	~complex64 | ~complex128
}

type Bool interface {
	~bool
	cmp.Ordered
}

type Pointer interface {
	~uintptr
}
