package threadsafe

import "sync"

var (
	// Any returns a new AnyObj instance
	Any = NewAnyObj
)

// AnyObj is a thread safe interface{}
type AnyObj struct {
	value interface{}
	l     sync.RWMutex
}

// NewAnyObj returns a new AnyObj
func NewAnyObj() *AnyObj {
	return &AnyObj{
		l: sync.RWMutex{},
	}
}

// Get returns thread safe []bytes
func (s *AnyObj) Get() interface{} {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.value
}

// Set sets []byte, thread safe
func (s *AnyObj) Set(value interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	s.value = value
}

