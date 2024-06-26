package utils

import (
	"sync"
	"time"
)

// TimeServer 是跟随服务器进程的值的时间，理论上其应该在服务器的每一帧逻辑帧进行一次更新
// TimeServer可以保证服务器运行时在每一帧的处理过程中所获取的时间都是一致的,
// 但是同样的相比于 time.Now() 其可能存在一定的延迟，该延迟最多不会超过一帧
// 如果在当前帧无法触发但是接近的时间任务 应该在服务器的下一帧执行

// 关于时间转换
// 1秒 Second = 1000毫秒 Millisecond
// 1毫秒 Millisecond = 1000微秒 Microsecond
// 1微秒 Microsecond = 1000纳秒 Nanosecond

var (
	TimeServerSingleton *timeServer
	serverOnce          sync.Once
)

func init() {
	serverOnce.Do(func() {
		TimeServerSingleton = newTimeServer()
	})
}

////// timeServer

type timeServer struct {
	now time.Time // 当前时间
}

func newTimeServer() *timeServer {
	return &timeServer{
		now: time.Now(),
	}
}

// Update TODO 务必在每一个逻辑帧进行更新
func (gs *timeServer) Update(now time.Time) {
	gs.now = now
}

func (gs *timeServer) Now() time.Time {
	return gs.now
}
