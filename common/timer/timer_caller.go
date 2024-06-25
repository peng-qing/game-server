package timer

import (
	"runtime/debug"

	"GameServer/gslog"
)

// ITimerCallback 定时器回调, 所有需要定时器的功能都应该实现该接口
type ITimerCallback interface {
	OnTimer(identifyID uint64, param any) bool
}

type ITimerCaller interface {
	ITimerCallback
	IdentifyID() uint64      // 获取定时器识别id
	CallbackParam() any      // 获取定时器回调参数
	CallTimerCallback() bool // 调用回调
	CallInterval() int64     // 调用间隔
	NextCallTime() int64     // 获取调用时间
}

type TimerCaller struct {
	ITimerCallback
	identifyID    uint64 // 定时器识别码
	callbackParam any    // 定时器回调参数
	nextCallTime  int64  // 下次调用时间
	callInterval  int64  // 调用间隔
}

func NewTimerCaller(identifyID uint64, callback ITimerCallback, param any, nextCallTime, callInterval int64) *TimerCaller {
	return &TimerCaller{
		ITimerCallback: callback,
		identifyID:     identifyID,
		callbackParam:  param,
		nextCallTime:   nextCallTime,
		callInterval:   callInterval,
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

func (gs *TimerCaller) CallInterval() int64 {
	if gs == nil {
		gslog.Error("[TimerCaller] GetCallInterval caller is nil")
		return 0
	}
	return gs.callInterval
}

func (gs *TimerCaller) NextCallTime() int64 {
	if gs == nil {
		gslog.Error("[TimerCaller] GetCallTime caller is nil")
		return 0
	}
	return gs.nextCallTime
}
