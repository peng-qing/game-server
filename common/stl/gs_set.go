package stl

import "GameServer/types"

// Set 集合
type Set[T comparable] struct {
	set map[T]types.None
}

// NewSet 返回一个集合
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		set: make(map[T]types.None),
	}
}

// Clear 清空所有元素
func (gs *Set[T]) Clear() {
	gs.set = make(map[T]types.None)
}

// IsEmpty 集合是否为空
func (gs *Set[T]) IsEmpty() bool {
	return len(gs.set) <= 0
}

// Contains 检查是否包含某个元素
func (gs *Set[T]) Contains(e T) bool {
	_, ok := gs.set[e]
	return ok
}

// Del 删除指定元素
func (gs *Set[T]) Del(e T) {
	delete(gs.set, e)
}

// Insert 插入元素
func (gs *Set[T]) Insert(e T) {
	gs.set[e] = types.None{}
}

// Size 集合大小
func (gs *Set[T]) Size() int {
	return len(gs.set)
}

// ToList 获取元素列表
func (gs *Set[T]) ToList() []T {
	list := make([]T, 0, gs.Size())
	for e := range gs.set {
		list = append(list, e)
	}
	return list
}
