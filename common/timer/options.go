package timer

type Options interface {
	apply(scheduler *TimeScheduler)
}

type OptionFunc func(scheduler *TimeScheduler)

func (f OptionFunc) apply(scheduler *TimeScheduler) {
	f(scheduler)
}
