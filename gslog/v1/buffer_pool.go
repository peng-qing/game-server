package v1

import "GameServer/common/pool"

// buffer pool 私有
// 只提供对应的Get方法来直接调用

var (
	_pool = pool.NewBufferPool()
	Get   = _pool.Get
)
