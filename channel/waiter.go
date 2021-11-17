package channel

import (
	"github.com/google/uuid"
	"go-threading/threadsafe"
	"sync"
)

type Waiter struct {
	id string
	c chan interface{}
	lock sync.Mutex
	waitCalled *threadsafe.SafeBool
}

func NewWaiter() *Waiter {
	return &Waiter{
		id: uuid.New().String(),
		c: make(chan interface{}),
		lock: sync.Mutex{},
		waitCalled: threadsafe.Bool(),
	}
}

// Wait will block until a new obj is passed or nil in case it got cancelled
func (w *Waiter) Wait() interface{} {
	w.waitCalled.Set(true)
	defer w.waitCalled.Set(false)
	obj := <- w.c
	return obj
}

// Fire will fire obj through the wait function if that function was called (and waiting), if not it will not fire.
func (w *Waiter) Fire(obj interface{}) {
	if w.waitCalled.Get() {
		w.c <- obj
	}
}