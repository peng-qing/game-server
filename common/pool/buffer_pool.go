package pool

import (
	"strconv"
	"time"
)

const (
	_bufferSize = 1024
)

// Buffer 缓冲区对象
// 不提供给NewBuffer方法
// 该对象的目的是为了池化缓冲区
type Buffer struct {
	buf  []byte
	pool *BufferPool
}

// AppendByte 写入byte
func (gs *Buffer) AppendByte(v byte) {
	gs.buf = append(gs.buf, v)
}

// AppendBytes 写入byte数组
func (gs *Buffer) AppendBytes(v []byte) {
	gs.buf = append(gs.buf, v...)
}

// AppendString 写入字符串
func (gs *Buffer) AppendString(s string) {
	gs.buf = append(gs.buf, s...)
}

// AppendInt 写入int 10进制
func (gs *Buffer) AppendInt(i int64) {
	gs.buf = strconv.AppendInt(gs.buf, i, 10)
}

// AppendUint 写入uint 10进制
func (gs *Buffer) AppendUint(i uint64) {
	gs.buf = strconv.AppendUint(gs.buf, i, 10)
}

// AppendTime 以指定格式写入时间字符串
func (gs *Buffer) AppendTime(t time.Time, layout string) {
	gs.buf = t.AppendFormat(gs.buf, layout)
}

// AppendFloat 写入float
func (gs *Buffer) AppendFloat(f float64, bitSize int) {
	gs.buf = strconv.AppendFloat(gs.buf, f, 'f', -1, bitSize)
}

// AppendBool 写入bool
func (gs *Buffer) AppendBool(v bool) {
	gs.buf = strconv.AppendBool(gs.buf, v)
}

// Size 缓冲区数据大小
func (gs *Buffer) Size() int {
	return len(gs.buf)
}

// Bytes 获取缓冲区数据 字节数组
func (gs *Buffer) Bytes() []byte {
	return gs.buf
}

// String 获取缓冲区数据 字符串
func (gs *Buffer) String() string {
	return string(gs.buf)
}

// Reset 清空缓冲区 重复利用底层数组避免频繁扩容
func (gs *Buffer) Reset() {
	gs.buf = gs.buf[:0]
}

// Write 实现 io.Writer
func (gs *Buffer) Write(bs []byte) (n int, err error) {
	gs.AppendBytes(bs)
	return len(bs), nil
}

// WriteByte 实现 io.ByteWriter
func (gs *Buffer) WriteByte(c byte) error {
	gs.AppendByte(c)
	return nil
}

// WriteString 实现 io.StringWriter
func (gs *Buffer) WriteString(s string) (n int, err error) {
	gs.AppendString(s)
	return len(s), nil
}

// Free 归还缓冲区对象给对象池
func (gs *Buffer) Free() {
	gs.pool.put(gs)
}

// TrimNewLine 去除缓冲区末尾的换行
func (gs *Buffer) TrimNewLine() {
	if i := len(gs.buf) - 1; i >= 0 && gs.buf[i] == '\n' {
		gs.buf = gs.buf[:i]
	}
}

type BufferPool struct {
	pool *Pool[*Buffer]
}

// NewBufferPool Buffer对象池
// 为什么不直接用泛型 Pool[*Buffer] 需要再包一层呢
// 主要是为了保证获取到的缓冲区对象都是干净的 不需要获取到后由上层再手段Reset
// 同时私有了Put方法 对象的放回由对象的Free方法完成
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: NewPool(func() *Buffer {
			return &Buffer{
				buf: make([]byte, 0, _bufferSize),
			}
		}),
	}
}

// Get 获取缓冲对象
func (gs *BufferPool) Get() *Buffer {
	buf := gs.pool.Get()
	buf.Reset()
	buf.pool = gs
	return buf
}

// put 私有不对外暴露
func (gs *BufferPool) put(buf *Buffer) {
	gs.pool.Put(buf)
}
