package timer

import (
	"GameServer/gslog"
	"context"
	"sync"
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

	return instance
}

func NewAutoTimeScheduler(options ...Options) *TimeScheduler {
	autoExecTimeScheduler := NewTimeScheduler(options...)

	return autoExecTimeScheduler
}

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
