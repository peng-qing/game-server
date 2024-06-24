package timer

type TimerManager struct {
	identifyIDToTimerCaller map[uint64]ITimerCaller
}

//func NewTimerManager() *TimerManager {
//
//}
