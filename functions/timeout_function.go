package functions

import (
	"time"
)

// NewTimeoutF will run a provided function in a new go routine for the provided timeout, after which it will be cancelled.
func NewTimeoutF(fn FuncWithStop, t time.Duration) *FuncResult {
	sf := NewStoppableF(fn)

	go func() {
		<-time.After(t)
		sf.Manager.Stop()
	}()

	res := sf.Start()
	return res
}
