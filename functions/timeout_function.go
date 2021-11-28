package functions

import (
	"go-threading/timer"
	"time"
)

// NewTimeoutF will run a provided function in a new go routine for the provided timeout, after which it will be cancelled.
func NewTimeoutF(fn FuncWithStop, t time.Duration) *FuncResult {
	sf := NewStoppableF(fn)

	tmr := timer.New()
	go func(tmr *timer.RoundTimer) {
		tmr.ResultChan().Wait()
		sf.Manager.Stop()
	}(tmr)
	tmr.Reset(t)

	res := sf.Start()
	tmr.Kill()
	return res
}
