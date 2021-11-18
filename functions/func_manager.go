package functions

import (
	"sync"
)

// FuncManager represents the object used to stop running functions
// should be used by the running function, once stopped the function act accordingly
type FuncManager interface {
	// IsStopped returns true if the funcManager already stopped
	IsStopped() bool
	// Stop will make IsStopped return true
	Stop()
}

type funcManager struct {
	stopped bool
	mut     sync.Mutex
}

func newFuncManager() *funcManager {
	s := funcManager{
		mut:    sync.Mutex{},
	}
	return &s
}

func (s *funcManager) IsStopped() bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.stopped
}

func (s *funcManager) Stop() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.stopped = true
}
