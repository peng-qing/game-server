package stl

// 对比 list.List 和 slice
// slice 的结构更简单, 维护成本低 性能大部分情况会由于list
// 如果元素过大，slice频繁拷贝则会耗费更多性能

const (
	initCapacity = 8
)

// Queue 队列
type Queue[T any] struct {
	data  []T
	begin int
	end   int
	cap   int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		data:  make([]T, initCapacity),
		begin: 0,
		end:   0,
		cap:   initCapacity,
	}
}

// Size 队列大小
func (gs *Queue[T]) Size() int {
	return gs.end - gs.begin
}

// Empty 队列是否为空
func (gs *Queue[T]) Empty() bool {
	return gs.Size() <= 0
}

// Clear 清空队列
func (gs *Queue[T]) Clear() {
	gs.data = make([]T, initCapacity)
	gs.begin = 0
	gs.end = 0
	gs.cap = initCapacity
}

// Push 往队尾提添加元素
func (gs *Queue[T]) Push(value T) {
	if gs.end >= gs.cap {
		// 扩容
		gs.expend()
	}
	gs.data[gs.end] = value
	gs.end++
}

func (gs *Queue[T]) expend() {
	// 优先首部前移
	if gs.begin > 0 {
		for i := 0; i < gs.end-gs.begin; i++ {
			gs.data[i] = gs.data[i+gs.begin]
		}
		gs.begin = 0
		gs.end -= gs.begin
		return
	}
	// 实际扩容
	if gs.cap < 1024 {
		// 翻倍
		gs.cap *= 2
	} else {
		gs.cap = gs.cap + (gs.cap / 4)
	}
	// 拷贝数据
	elements := make([]T, gs.cap)
	copy(elements, gs.data[gs.begin:gs.end])
	gs.data = elements
}

// Pop 出队 通过具名变量处理返回空值
func (gs *Queue[T]) Pop() (element T) {
	if gs.Empty() {
		return
	}

	element = gs.data[gs.begin]
	gs.begin++

	return
}
