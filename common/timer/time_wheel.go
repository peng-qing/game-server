package timer

import (
	"GameServer/gslog"
	"context"
	"sync"
	"time"
)

// TODO 多级时间轮不会直接触发回调

// TimeWheel 多级时间轮
type TimeWheel struct {
	name          string                         // 时间轮标识
	interval      time.Duration                  // 刻度时间间隔 单位ms
	scales        int                            // 刻度数
	current       int                            // 当前刻度
	timeQueue     map[int]map[int64]ITimerCaller // 时间轮上所有的Timer identifyID => TimerCaller
	onceStop      sync.Once                      // 保证只进行一次关闭
	nextTimeWheel *TimeWheel                     // 下一级时间轮
	ctx           context.Context                // context 用于关闭时间轮
	ticker        *time.Ticker                   // 用于时间轮转动
	mutex         sync.Mutex                     // 锁
}

///// constructors

// NewTimeWheel 创建多级时间轮
func NewTimeWheel(ctx context.Context, name string, interval time.Duration, scales int) *TimeWheel {
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
		name:          name,
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

	callTime := caller.NextCallTime()
	if callTime.IsZero() {
		gslog.Error("[TimeWheel] caller time is zero", "identifyID", identifyID)
		return
	}
	if callTime.Before(time.Now()) {
		gslog.Warn("[TimeWheel] caller time before now", "callTime", callTime.Unix(), "identifyID", identifyID)
		return
	}
	// 需要跨越几个刻度
	delayScales := callTime.Sub(time.Now()) / gs.interval
	// 如果大于一个刻度
	if delayScales >= 1 {
		targetScales := (gs.current + int(delayScales)) % gs.scales
		gs.timeQueue[targetScales][identifyID] = caller
		// 这里不处理后续刻度的定时器是否在最底层时间轮
		// 在每次轮转的时候 重新添加相关刻度的所有时间轮
		return
	}
	// 如果没有下一级时间轮
	if gs.nextTimeWheel == nil {
		// 当前是底层时间轮
		if forceNext {
			// 强制移入下一刻度 这部分发生在时间轮自转
			// 如果当前刻度不未执行且不移入下一刻度 其会失去被调度的计划
			// 将其强制移入下一个时间轮刻度 等待下一个轮转前被调度
			gs.timeQueue[(gs.current+1)%gs.scales][identifyID] = caller
			return
		}
		// 手动添加
		gs.timeQueue[gs.current][identifyID] = caller
		return
	}
	// 下一级时间轮
	gs.nextTimeWheel.AddTimer(identifyID, caller)
}

// run run在子goroutine中处理 所以相关接口需要单独加锁
func (gs *TimeWheel) run() {
	defer func() {
		if err := recover(); err != nil {
			gslog.Critical("[TimeWheel] run panic..", "err", err)
		}
	}()

	for {
		select {
		case <-gs.ticker.C:
			// 到时间 转刻度
			gs.mutex.Lock()
			gs.rotate()
			gs.mutex.Unlock()
		case <-gs.ctx.Done():
			// 时间轮关闭 退出
			gs.stop()
			return
		}
	}
}

func (gs *TimeWheel) rotate() {
	// 取出当前刻度全部定时器
	curScaleTimers := gs.timeQueue[gs.current]
	gs.timeQueue[gs.current] = make(map[int64]ITimerCaller)
	for identifyID, caller := range curScaleTimers {
		// 自转
		gs.addTimer(identifyID, caller, true)
	}
	// 下一个刻度的定时器也重新添加
	nextScaleTimers := gs.timeQueue[(gs.current+1)%gs.scales]
	gs.timeQueue[(gs.current+1)%gs.scales] = nextScaleTimers
	for identifyID, caller := range nextScaleTimers {
		gs.addTimer(identifyID, caller, true)
	}
	// 刻度移动
	gs.current = (gs.current + 1) % gs.scales
}

func (gs *TimeWheel) stop() {
	gs.onceStop.Do(func() {
		gs.ticker.Stop()
	})
}

////// exports

// AddTimer 添加定时器
func (gs *TimeWheel) AddTimer(identifyID int64, caller ITimerCaller) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.addTimer(identifyID, caller, false)
}

// RemoveTimer 移除定时器
func (gs *TimeWheel) RemoveTimer(identifyID int64) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	for i := 0; i < gs.scales; i++ {
		if _, ok := gs.timeQueue[i][identifyID]; ok {
			delete(gs.timeQueue[i], identifyID)
		}
	}
}

// AddTimerWheel 添加下一级时间轮
func (gs *TimeWheel) AddTimerWheel(nextTimerWheel *TimeWheel) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.nextTimeWheel = nextTimerWheel
}

// GetTimerWithDuration 获取一定时间间隔内的所有timer
func (gs *TimeWheel) GetTimerWithDuration(interval time.Duration) map[int64]ITimerCaller {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// 找到最底层时间轮
	curTimeWheel := gs
	if curTimeWheel.nextTimeWheel != nil {
		curTimeWheel = curTimeWheel.nextTimeWheel
	}

	now := time.Now()
	timerMap := make(map[int64]ITimerCaller)

	for identifyID, caller := range curTimeWheel.timeQueue[curTimeWheel.current] {
		callTime := caller.NextCallTime()
		if callTime.Before(now.Add(interval)) {
			// 定时器已经超时
			timerMap[identifyID] = caller
			delete(curTimeWheel.timeQueue[curTimeWheel.current], identifyID)
		}
	}

	return timerMap
}

func (gs *TimeWheel) Start() {
	gs.ticker = time.NewTicker(gs.interval)

	go gs.run()
	gslog.Info("[TimeWheel] time wheel start....", "interval", gs.interval, "scales", gs.scales)
}
