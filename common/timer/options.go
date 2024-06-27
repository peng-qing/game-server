package timer

type Options interface {
	apply(scheduler *TimerScheduler)
}

type OptionFunc func(scheduler *TimerScheduler)

func (f OptionFunc) apply(scheduler *TimerScheduler) {
	f(scheduler)
}
