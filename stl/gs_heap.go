package stl

import (
	"cmp"
	"container/heap"
	"fmt"
)

// HeapVal go heap 元素
type HeapVal[K cmp.Ordered, V any] struct {
	Key   K      // 用于比较的Key
	Val   V      // 参数
	UUID  uint64 // 唯一标识 Universally Unique Identifier
	GUID  uint64 // 唯一标识 Globally Unique Identifier
	index int    // 位置索引
}

// Heap 小根堆
type Heap[K cmp.Ordered, V any] struct {
	heap *gsHeap[K, V]
}

// NewHeap 创建一个小根堆
func NewHeap[K cmp.Ordered, V any]() *Heap[K, V] {
	gs := &Heap[K, V]{
		heap: &gsHeap[K, V]{
			val: make([]*HeapVal[K, V], 0),
		},
	}
	//heap.Init(gs.heap)
	return gs
}

// Push 往小根堆中添加元素
func (gs *Heap[K, V]) Push(x *HeapVal[K, V]) {
	heap.Push(gs.heap, x)
}

// Pop 获取最后一个元素
func (gs *Heap[K, V]) Pop() *HeapVal[K, V] {
	x := heap.Pop(gs.heap)
	if res, ok := x.(*HeapVal[K, V]); ok {
		return res
	}
	return nil
}

// Len 小根堆元素长度
func (gs *Heap[K, V]) Len() int {
	return gs.heap.Len()
}

// MinHeapVal 获取最小的元素
func (gs *Heap[K, V]) MinHeapVal() *HeapVal[K, V] {
	if gs.Len() > 0 {
		return gs.heap.val[0]
	}
	return nil
}

// RemoveByUUID 根据 UUID 删除堆元素
func (gs *Heap[K, V]) RemoveByUUID(uuid uint64) bool {
	size := gs.Len()
	for i := 0; i < size; i++ {
		val := gs.heap.val[i]
		if val == nil {
			return false
		}
		if val.UUID == uuid {
			heap.Remove(gs.heap, i)
			return true
		}
	}
	return false
}

// RemoveByGUID 根据 GUID 删除堆元素
func (gs *Heap[K, V]) RemoveByGUID(guid uint64) bool {
	size := gs.Len()
	for i := 0; i < size; i++ {
		val := gs.heap.val[i]
		if val == nil {
			return false
		}
		if val.GUID == guid {
			heap.Remove(gs.heap, i)
			return true
		}
	}
	return false
}

// RemoveByGUIDAndUUID 根据 UUID 和 GUID 删除堆元素
func (gs *Heap[K, V]) RemoveByGUIDAndUUID(guid, uuid uint64) bool {
	size := gs.Len()
	for i := 0; i < size; i++ {
		val := gs.heap.val[i]
		if val == nil {
			return false
		}
		if val.GUID == guid && val.UUID == uuid {
			heap.Remove(gs.heap, i)
			return true
		}
	}
	return false
}

// RemoveAllByUUID 根据UUID删除所有堆元素
// @returns 移除的元素个数
func (gs *Heap[K, V]) RemoveAllByUUID(uuid uint64) int {
	removeCount := 0
	size := gs.Len()
	for i := 0; i < size; i++ {
		val := gs.heap.val[i]
		if val == nil {
			return removeCount
		}
		if val.UUID == uuid {
			heap.Remove(gs.heap, i)
			removeCount++
		}
	}
	return removeCount
}

// RemoveAllByGUID 根据GUID删除所有堆元素
func (gs *Heap[K, V]) RemoveAllByGUID(guid uint64) int {
	removeCount := 0
	size := gs.Len()
	for i := 0; i < size; i++ {
		val := gs.heap.val[i]
		if val == nil {
			return removeCount
		}
		if val.GUID == guid {
			heap.Remove(gs.heap, i)
			removeCount++
		}
	}
	return removeCount
}

func (gs *Heap[K, V]) Debug() {
	for i := 0; i < gs.Len(); i++ {
		fmt.Printf("HeapVal: Key=%+v, Val=%+v, UUID=%d, GUID=%d, index=%d, i=%d\n",
			gs.heap.val[i].Key, gs.heap.val[i].Val, gs.heap.val[i].UUID, gs.heap.val[i].GUID, gs.heap.val[i].index, i)
	}
}

// =========================== 实现go的Heap接口 =========================== //

// 此处实现的是小根堆 pop出的顺序为小到大
type gsHeap[K cmp.Ordered, V any] struct {
	val []*HeapVal[K, V]
}

func (gs *gsHeap[K, V]) Len() int {
	return len(gs.val)
}

// Less 比较 < 实现的是小根堆 > 实现的是大根堆
func (gs *gsHeap[K, V]) Less(i, j int) bool {
	return gs.val[i].Key < gs.val[j].Key
}

func (gs *gsHeap[K, V]) Swap(i, j int) {
	gs.val[i], gs.val[j] = gs.val[j], gs.val[i]
	gs.val[i].index = i
	gs.val[j].index = j
}

func (gs *gsHeap[K, V]) Push(x any) {
	elem, ok := x.(*HeapVal[K, V])
	if !ok {
		fmt.Printf("[gsHeap] Push type error, x:%+v", x)
		return
	}
	elem.index = len(gs.val)
	gs.val = append(gs.val, elem)
}

func (gs *gsHeap[K, V]) Pop() any {
	n := len(gs.val)
	elem := gs.val[n-1]
	elem.index = -1 // 安全起见标记为无效元素

	gs.val[n-1] = nil // 避免内存泄露 avoid memory leak
	gs.val = gs.val[0 : n-1]

	return elem
}
