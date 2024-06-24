package timer

import (
	"GameServer/gslog"
	"runtime/debug"
)

// ITimerCallback 定时器回调, 所有需要定时器的功能都应该实现该接口
type ITimerCallback interface {
	OnTimer(identifyID uint64, param any) bool
}

type ITimerCaller interface {
	ITimerCallback
	IdentifyID() uint64                       // 获取定时器识别id
	CallbackParam() any                       // 获取定时器回调参数
	CallTimerCallback() bool                  // 调用回调
	GetCallInterval() int64                   // 调用间隔
	GetNextCallTime() int64                   // 获取调用时间
	TryUpdateNextCallTime(curTime int64) bool // 尝试更新下次调用时间
}

type TimerCaller struct {
	ITimerCallback
	identifyID    uint64 // 定时器识别码
	callbackParam any    // 定时器回调参数
	nextCallTime  int64  // 下次调用时间
	callInterval  int64  // 调用间隔
	lastCallTime  int64  // 最后一次触发时间
}

func NewTimerCaller(identifyID uint64, callback ITimerCallback, param any, nextCallTime, callInterval, lastCallTime int64) *TimerCaller {
	return &TimerCaller{
		ITimerCallback: callback,
		identifyID:     identifyID,
		callbackParam:  param,
		nextCallTime:   nextCallTime,
		callInterval:   callInterval,
		lastCallTime:   lastCallTime,
	}
}

func (gs *TimerCaller) IdentifyID() uint64 {
	if gs == nil {
		gslog.Error("[TimerCaller] IdentifyID caller is nil")
		return 0
	}
	return gs.identifyID
}

func (gs *TimerCaller) CallbackParam() any {
	if gs == nil {
		gslog.Error("[TimerCaller] CallbackParam caller is nil")
		return nil
	}
	return gs.callbackParam
}

func (gs *TimerCaller) CallTimerCallback() bool {
	if gs == nil {
		gslog.Error("[TimerCaller] CallTimerCallback caller is nil")
		return false
	}
	defer func() {
		if err := recover(); err != nil {
			gslog.Critical("[TimerCaller] CallTimerCallback  timer callback panic", gslog.Any("err", err), gslog.String("stack", debug.Stack()))
			return
		}
	}()

	return gs.OnTimer(gs.IdentifyID(), gs.CallbackParam())
}

func (gs *TimerCaller) GetCallInterval() int64 {
	if gs == nil {
		gslog.Error("[TimerCaller] GetCallInterval caller is nil")
		return 0
	}
	return gs.callInterval
}

func (gs *TimerCaller) GetNextCallTime() int64 {
	if gs == nil {
		gslog.Error("[TimerCaller] GetCallTime caller is nil")
		return 0
	}
	return gs.nextCallTime
}

// TryUpdateNextCallTime 尝试更新下次调用时间
// @param curTime 当前时间
// @param lastCallTime 最后一次调用时间
func (gs *TimerCaller) TryUpdateNextCallTime(curTime int64) bool {
	if gs == nil {
		gslog.Error("[TimerCaller] TryUpdateNextCallTime caller is nil")
		return false
	}
	gslog.Trace("[TimerCaller] TryUpdateNextCallTime start", "IdentifyID", gs.IdentifyID(), "curTime", curTime)
	callTime := gs.GetNextCallTime()
	callInterval := gs.GetCallInterval()
	lastCallTime := gs.LastCallTime()

	if lastCallTime <= 0 && callInterval <= 0 {
		// 一次性定时器
		if callTime > curTime {
			// 还没触发
			gslog.Trace("[TimerCaller] TryUpdateNextCallTime timer not expired", "IdentifyID", gs.IdentifyID(),
				"curTime", curTime, "callTime", callTime, "lastCallTime", lastCallTime)
			return false
		}
		// 过期了
		gslog.Debug("[TimerCaller] TryUpdateNextCallTime timer expired", "IdentifyID", gs.IdentifyID(),
			"curTime", curTime, "callTime", callTime, "lastCallTime", lastCallTime)
		gs.nextCallTime = 0
	}
	if lastCallTime > 0 {
		// 超时失效
		if callTime < curTime {
			gslog.Warn("[TimerCaller] TryUpdateNextCallTime timer expired", "IdentifyID", gs.IdentifyID(),
				"curTime", curTime, "callTime", callTime, "lastCallTime", lastCallTime)
			return false
		}
		// 下次触发在上次触发之前(改时间导致异常)
		if callTime <= lastCallTime {
			gslog.Warn("[TimerCaller] TryUpdateNextCallTime time modified", "IdentifyID", gs.IdentifyID(),
				"curTime", curTime, "callTime", callTime, "lastCallTime", lastCallTime)
			return false
		}
		// 不需要调整
		return false
	}

	return true
}

func (gs *TimerCaller) LastCallTime() int64 {
	if gs == nil {
		gslog.Error("[TimerCaller] LastCallTime caller is nil")
		return 0
	}
	return gs.lastCallTime
}
