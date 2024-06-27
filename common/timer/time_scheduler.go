package timer

import (
	"GameServer/gslog"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type TimeScheduler struct {
	IdentifyID   int64              // 定时器自增标识ID
	maxDelay     time.Duration      // 最大延迟误差时间
	chanSize     int                // 执行缓存队列大小
	topTimeWheel *TimeWheel         // 顶层时间轮
	triggerChan  chan ITimerCaller  // 执行队列
	ticker       *time.Ticker       // 执行定时器
	ctx          context.Context    // context
	cancel       context.CancelFunc // 关闭函数
	mutex        sync.Mutex         // 互斥锁
}

// NewTimeScheduler 创建时间轮调度器
func NewTimeScheduler(options ...Options) *TimeScheduler {
	instance := &TimeScheduler{
		IdentifyID: 0,
		maxDelay:   defaultMaxDelayDuration,
		chanSize:   defaultMaxCallChanSize,
	}

	// 创建多级时间轮
	ctx, cancel := context.WithCancel(context.Background())
	instance.ctx = ctx
	instance.cancel = cancel

	for _, option := range options {
		option.apply(instance)
	}

	instance.triggerChan = make(chan ITimerCaller, instance.chanSize)
	instance.ticker = time.NewTicker(instance.maxDelay / 2)

	if instance.topTimeWheel == nil {
		hourWheel := NewTimeWheel(ctx, HourWheelName, HourInterval, HourScales)
		minuteWheel := NewTimeWheel(ctx, MinuteWheelName, MinuteInterval, MinuteScales)
		secondWheel := NewTimeWheel(ctx, SecondWheelName, SecondInterval, SecondScales)

		hourWheel.AddTimerWheel(minuteWheel)
		minuteWheel.AddTimerWheel(secondWheel)

		instance.topTimeWheel = hourWheel
	}

	instance.Run()

	return instance
}

// NewAutoTimeScheduler 创建时间轮调度器 自动任务调度
func NewAutoTimeScheduler(options ...Options) *TimeScheduler {
	autoExecTimeScheduler := NewTimeScheduler(options...)

	autoExecTimeScheduler.Execute()

	return autoExecTimeScheduler
}

// Run 后台线程处理调度
func (gs *TimeScheduler) Run() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				gslog.Critical("[TimeScheduler] Run scheduler run panic..", "err", err)
			}
		}()
		for {
			select {
			case <-gs.ctx.Done():
				gs.ticker.Stop()
				return
			case <-gs.ticker.C:
				now := time.Now()
				timerMap := gs.topTimeWheel.GetTimerWithDuration(gs.maxDelay)
				for identifyID, caller := range timerMap {
					callTime := caller.NextCallTime()
					if now.Sub(callTime) > gs.maxDelay {
						// 超时
						gslog.Warn("[TimeScheduler] Run scheduler run time exceed max delay",
							"callTime", callTime.Unix(), "now", now.Unix(), "identifyID", identifyID)
					}
					gs.triggerChan <- caller
				}
			}
		}
	}()
}

// Execute 自动执行调度时调用
func (gs *TimeScheduler) Execute() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				gslog.Critical("[TimeScheduler] Execute scheduler exec panic..", "err", err)
			}
		}()
		for {
			select {
			case <-gs.ctx.Done():
				return
			case caller := <-gs.triggerChan:
				caller.CallTimerCallback()
			}
		}
	}()
}

// TriggerChan 获取任务执行队列
func (gs *TimeScheduler) TriggerChan() chan ITimerCaller {
	return gs.triggerChan
}

// Stop 停止调度 多级时间轮同时停止
func (gs *TimeScheduler) Stop() {
	gs.cancel()
	gs.ticker.Stop()
	close(gs.triggerChan)
}

// AddTimer 添加定时器
func (gs *TimeScheduler) AddTimer(callback ITimerCallback, param any, nextCallTime time.Time, callInterval time.Duration) int64 {
	if gs == nil {
		gslog.Error("[TimeScheduler] AddTimer called but TimeScheduler is nil")
		return 0
	}

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	identifyID := atomic.LoadInt64(&gs.IdentifyID)
	timerCaller := NewTimerCaller(identifyID, callback, param, nextCallTime, callInterval)

	gs.topTimeWheel.AddTimer(identifyID, timerCaller)

	return identifyID
}

// CancelTimer 关闭注册定时器
func (gs *TimeScheduler) CancelTimer(identifyID int64) {
	if gs == nil {
		gslog.Error("[TimeScheduler] CancelTimer called but TimeScheduler is nil")
		return
	}

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	curTimeWheel := gs.topTimeWheel

	for curTimeWheel != nil {
		curTimeWheel.RemoveTimer(identifyID)
		curTimeWheel = curTimeWheel.nextTimeWheel
	}
}
