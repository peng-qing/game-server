package timer

import "time"

const (
	// 默认时间轮刻度数
	defaultTimeWheelScales = 12
	// 默认最大执行回调缓冲队列大小
	defaultMaxCallChanSize = 2048
	// 默认最大误差时间
	defaultMaxDelayDuration = 100 * time.Millisecond
)

// 一些默认的时间轮配置
const (
	HourWheelName = "Hour"
	HourScales    = 24
	HourInterval  = time.Hour

	MinuteWheelName = "Minute"
	MinuteScales    = 60
	MinuteInterval  = time.Minute

	SecondWheelName = "Second"
	SecondScales    = 60
	SecondInterval  = time.Second
)
