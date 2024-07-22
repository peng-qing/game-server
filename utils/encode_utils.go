package utils

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidBuffer = errors.New("invalid bytes buffer")
)

// EncodeVariableInt 变长int/int64编码 对于小于128的值采用单字节编码
// 高出的值采取 低7位编码有效数据,可以编码128个数值; 最高位延续位指示后续是否还有剩余字节
func EncodeVariableInt(num int64) []byte {
	var enc []byte

	for {
		digit := byte(num % 128)
		num /= 128
		if num > 0 {
			digit |= 0x80
		}
		enc = append(enc, digit)
		if num == 0 {
			break
		}
	}

	return enc
}

// DecodeReaderVariableInt  变长int解码
func DecodeReaderVariableInt(r io.Reader) (int, error) {
	var num uint32
	var shift uint32
	var err error
	var bs = make([]byte, 1)

	for shift < 28 {
		_, err = r.Read(bs)
		if err != nil {
			return 0, err
		}
		digit := bs[0]

		num |= uint32(digit&0x7F) << shift
		if digit&0x80 == 0 {
			break
		}
		shift += 7
	}

	return int(num), err
}

// DecodeVariableInt  变长int解码
func DecodeVariableInt(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, ErrInvalidBuffer
	}

	var num uint32
	var shift uint32

	for shift < 28 {
		if len(data) <= 0 {
			break
		}
		digit := data[0]
		data = data[1:]

		num |= uint32(digit&0x7F) << shift
		if digit&0x80 == 0 {
			break
		}
		shift += 7
	}

	return int(num), nil
}

func DecodeVariableInt64(data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, ErrInvalidBuffer
	}

	var num uint64
	var shift uint64

	for shift < 56 {
		if len(data) <= 0 {
			break
		}
		digit := data[0]
		data = data[1:]

		num |= uint64(digit&0x7F) << shift
		if digit&0x80 == 0 {
			break
		}
		shift += 7
	}

	return int64(num), nil
}

func DecodeReaderVariableInt64(r io.Reader) (int64, error) {
	var num uint64
	var shift uint64
	var err error
	var bs = make([]byte, 1)

	for shift < 56 {
		_, err = r.Read(bs)
		if err != nil {
			return 0, err
		}
		digit := bs[0]

		num |= uint64(digit&0x7F) << shift
		if digit&0x80 == 0 {
			break
		}
		shift += 7
	}

	return int64(num), err
}

func EncodeBytes(field []byte, order binary.ByteOrder) []byte {
	length := make([]byte, 4)
	order.PutUint32(length, uint32(len(field)))
	return append(length, field...)
}

func DecodeBytes(data []byte, order binary.ByteOrder) ([]byte, error) {
	if len(data) < 4 {
		return nil, ErrInvalidBuffer
	}
	length := order.Uint32(data[:4])
	field := make([]byte, length)
	copy(field, data[4:])

	return field, nil
}

func DecodeReaderBytes(r io.Reader, order binary.ByteOrder) ([]byte, error) {
	num := make([]byte, 4)
	_, err := r.Read(num)
	if err != nil {
		return nil, err
	}
	length := order.Uint32(num)

	field := make([]byte, length)
	_, err = r.Read(field)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func EncodeString(str string, order binary.ByteOrder) []byte {
	return EncodeBytes([]byte(str), order)
}

func DecodeString(data []byte, order binary.ByteOrder) (string, error) {
	buf, err := DecodeBytes(data, order)
	return string(buf), err
}

func DecodeReaderString(r io.Reader, order binary.ByteOrder) (string, error) {
	buf, err := DecodeReaderBytes(r, order)
	return string(buf), err
}
