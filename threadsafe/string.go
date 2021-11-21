package threadsafe

import "sync"

var (
	// String returns a new SafeString instance
	String = NewString
)

// SafeString is a thread safe string
type SafeString struct {
	value string
	l     sync.RWMutex
}

// NewString returns a new String
func NewString(s string) *SafeString {
	return &SafeString{
		l:     sync.RWMutex{},
		value: s,
	}
}

// Get returns thread safe string
func (s *SafeString) Get() string {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.value
}

// Set sets string, thread safe
func (s *SafeString) Set(value string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.value = value
}
