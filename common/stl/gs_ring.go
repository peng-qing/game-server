package stl

type RingQueue[T any] struct {
	buffer []T
	head   int
	tail   int
	size   int
	length int
}

func NewRingQueue[T any](size int) *RingQueue[T] {
	return &RingQueue[T]{
		buffer: make([]T, size),
		head:   0,
		tail:   0,
		length: 0,
		size:   size,
	}
}

func (gs *RingQueue[T]) Push(value T) {
	if gs.tail == gs.head && !gs.Empty() {
		gs.expend()
	}
	gs.buffer[gs.tail] = value
	gs.tail = (gs.tail + 1) % gs.size
	gs.length++
}

func (gs *RingQueue[T]) expend() {
	oldSize := gs.size
	if gs.size <= 1024 {
		gs.size *= 2
	} else {
		gs.size += gs.size / 4
	}

	newBuffer := make([]T, gs.size)
	for i := 0; i < oldSize; i++ {
		newBuffer[i] = gs.buffer[(gs.tail+i)%oldSize]
	}
	gs.head = 0
	gs.tail = gs.length

	gs.buffer = newBuffer
}

func (gs *RingQueue[T]) Empty() bool {
	return gs.length == 0
}

func (gs *RingQueue[T]) Length() int {
	return gs.length
}

func (gs *RingQueue[T]) Pop() (res T) {
	if gs.Empty() {
		return res
	}
	var e T
	res = gs.buffer[gs.head]
	gs.buffer[gs.head] = e
	gs.head = (gs.head + 1) % gs.size
	gs.length--

	return res
}

func (gs *RingQueue[T]) PopMany(length int) []T {
	if length <= 0 || gs.Empty() {
		return nil
	}
	if length > gs.length {
		length = gs.length
	}
	gs.length -= length
	var e T

	res := make([]T, length)
	for i := 0; i < length; i++ {
		res[i] = gs.buffer[gs.head+i%gs.size]
		gs.buffer[gs.head+i%gs.size] = e
	}
	gs.head = (gs.head + length) % gs.size

	return res
}
