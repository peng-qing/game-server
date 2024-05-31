package pool

import "sync"

// 通用的对象池

type Pool[T any] struct {
	pool sync.Pool
}

func NewPool[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}

func (gs *Pool[T]) Get() T {
	return gs.pool.Get().(T)
}

func (gs *Pool[T]) Put(element T) {
	gs.pool.Put(element)
}
