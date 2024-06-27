package timer

import (
	"runtime/debug"
	"time"

	"GameServer/gslog"
)

// ITimerCallback 定时器回调, 所有需要定时器的功能都应该实现该接口
type ITimerCallback interface {
	OnTimer(identifyID int64, param any) bool
}

type ITimerCaller interface {
	ITimerCallback
	IdentifyID() int64           // 获取定时器识别id
	CallbackParam() any          // 获取定时器回调参数
	CallTimerCallback() bool     // 调用回调
	CallInterval() time.Duration // 调用间隔
	NextCallTime() time.Time     // 获取调用时间
}

type TimerCaller struct {
	ITimerCallback
	identifyID    int64         // 定时器识别码
	callbackParam any           // 定时器回调参数
	nextCallTime  time.Time     // 下次调用时间
	callInterval  time.Duration // 调用间隔
}

func NewTimerCaller(identifyID int64, callback ITimerCallback, param any, nextCallTime time.Time, callInterval time.Duration) *TimerCaller {
	return &TimerCaller{
		ITimerCallback: callback,
		identifyID:     identifyID,
		callbackParam:  param,
		nextCallTime:   nextCallTime,
		callInterval:   callInterval,
	}
}

func (gs *TimerCaller) IdentifyID() int64 {
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

func (gs *TimerCaller) CallInterval() time.Duration {
	if gs == nil {
		gslog.Error("[TimerCaller] GetCallInterval caller is nil")
		return 0
	}
	return gs.callInterval
}

func (gs *TimerCaller) NextCallTime() time.Time {
	if gs == nil {
		gslog.Error("[TimerCaller] GetCallTime caller is nil")
		return time.Time{}
	}

	return gs.nextCallTime
}
