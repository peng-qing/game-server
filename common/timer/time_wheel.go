package timer

import (
	"GameServer/gslog"
	"context"
	"sync"
	"time"
)

// TODO 多级时间轮不会创建异步后台goroutine去创建/检索/添加Timer等操作

// TimeWheel 多级时间轮
type TimeWheel struct {
	interval      time.Duration                  // 刻度时间间隔 单位ms
	scales        int                            // 刻度数
	current       int                            // 当前刻度
	timeQueue     map[int]map[int64]ITimerCaller // 时间轮上所有的Timer identifyID => TimerCaller
	onceStop      sync.Once                      // 保证只进行一次关闭
	nextTimeWheel *TimeWheel                     // 下一级时间轮
	ctx           context.Context                // context 用于关闭时间轮
	mutex         sync.Mutex                     // 锁
}

///// constructors

// NewTimeWheel 创建多级时间轮
func NewTimeWheel(ctx context.Context, interval time.Duration, scales int) *TimeWheel {
	if ctx == nil {
		ctx = context.Background()
	}
	if interval == 0 {
		interval = time.Second
	}
	if scales == 0 {
		scales = defaultTimeWheelScales
	}

	instance := &TimeWheel{
		interval:      interval,
		scales:        scales,
		current:       0,
		timeQueue:     make(map[int]map[int64]ITimerCaller),
		onceStop:      sync.Once{},
		nextTimeWheel: nil,
		ctx:           ctx,
		mutex:         sync.Mutex{},
	}

	for i := 0; i < scales; i++ {
		instance.timeQueue[i] = make(map[int64]ITimerCaller)
	}

	return instance
}

//////// internal

// addTimer 将定时器添加到多级时间轮
// @param identifyID 定时器唯一标识
// @param caller 回调对象
// @param forceNext 强制添加到下一级时间轮
func (gs *TimeWheel) addTimer(identifyID int64, caller ITimerCaller, forceNext bool) {
	defer func() {
		if err := recover(); err != nil {
			gslog.Error("[TimeWheel] addTimer err", "error", err)
			return
		}
	}()

}

////// exports

func (gs *TimeWheel) AddTimer(identifyID int64, caller ITimerCaller) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.addTimer(identifyID, caller, false)
}
