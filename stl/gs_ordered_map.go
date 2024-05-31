package stl

import "cmp"

type OrderedMapElement[K cmp.Ordered] struct {
	key  K
	next *OrderedMapElement[K]
}

func NewOrderedMapElement[K cmp.Ordered](key K) *OrderedMapElement[K] {
	return &OrderedMapElement[K]{
		key:  key,
		next: nil,
	}
}

// Key 有序Map元素排序Key
func (gs *OrderedMapElement[K]) Key() K {
	return gs.key
}

// Next 下一个元素
func (gs *OrderedMapElement[K]) Next() *OrderedMapElement[K] {
	return gs.next
}

// OrderedMap 有序Map
type OrderedMap[K cmp.Ordered, V any] struct {
	pHead      *OrderedMapElement[K]
	mapElement map[K]V
}

// NewOrderedMap 创建一个有序map
func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		pHead:      nil,
		mapElement: make(map[K]V),
	}
}

// Size 有序Map大小
func (gs *OrderedMap[K, V]) Size() int {
	return len(gs.mapElement)
}

// Clear 清空有序map所有元素
func (gs *OrderedMap[K, V]) Clear() {
	gs.pHead = nil
	gs.mapElement = make(map[K]V)
}

// MapElements 获取所有Map元素
func (gs *OrderedMap[K, V]) MapElements() map[K]V {
	return gs.mapElement
}

// Begin 获取有序Map头元素
func (gs *OrderedMap[K, V]) Begin() *OrderedMapElement[K] {
	return gs.pHead
}

// Delete 删除元素
func (gs *OrderedMap[K, V]) Delete(key K) {
	_, exist := gs.mapElement[key]
	if !exist {
		return
	}
	delete(gs.mapElement, key)
	if gs.pHead == nil {
		return
	}
	if gs.pHead.Key() == key {
		gs.pHead = gs.pHead.Next()
		return
	}
	cursor := gs.pHead
	for i := 0; i < gs.Size(); i++ {
		if cursor.Next() == nil {
			break
		}
		if cursor.Next().Key() == key {
			cursor.next = cursor.Next().Next()
			break
		}
		cursor = cursor.Next()
	}
}

// Contains 指定Key是否存在
func (gs *OrderedMap[K, V]) Contains(key K) bool {
	_, exist := gs.mapElement[key]
	return exist
}

// Insert 插入键值对
func (gs *OrderedMap[K, V]) Insert(key K, val V) {
	_, exist := gs.mapElement[key]
	gs.mapElement[key] = val
	if exist {
		return
	}
	cursor := gs.pHead
	element := NewOrderedMapElement(key)
	if cursor == nil {
		gs.pHead = element
		return
	}
	if cursor.Key() >= element.Key() {
		element.next = cursor
		gs.pHead = element
		return
	}
	for i := 0; i < gs.Size(); i++ {
		if cursor.Next() == nil {
			cursor.next = element
			break
		}
		if cursor.Next().Key() >= element.Key() {
			element.next = cursor.next
			cursor.next = element
			break
		}
		cursor = cursor.next
	}
}

// Get 获取元素
func (gs *OrderedMap[K, V]) Get(key K) V {
	return gs.mapElement[key]
}
