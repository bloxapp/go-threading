package channel

import (
	"context"
	"go-threading/queue"
	"go-threading/threadsafe"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const QueueSize = 5

type Waiter struct {
	id         string
	c          chan interface{}
	lock       sync.Mutex
	waitCalled *threadsafe.SafeBool
	queue      queue.Queue
}

func NewWaiter() *Waiter {
	return &Waiter{
		id:         uuid.New().String(),
		c:          make(chan interface{}),
		lock:       sync.Mutex{},
		waitCalled: threadsafe.Bool(),
		queue:      queue.New(queue.FIFO, QueueSize),
	}
}

// Wait will block until a new obj is passed from a queue or from firing.
// If queue has items, will return immediately after popping firs item
func (w *Waiter) Wait() interface{} {
	w.waitCalled.Set(true)
	defer w.waitCalled.Set(false)

	// check queue and return
	if w.queue.Len() > 0 {
		return w.queue.Pop()
	}

	// no queue, wait
	obj := <-w.c
	return obj
}

// WaitWithTimeout will return a fired object or an error if deadline exceeded
func (w *Waiter) WaitWithTimeout(duration time.Duration) interface{} {
	c, _ := context.WithTimeout(context.Background(), duration)
	return w.WaitWithContext(c)
}

// WaitWithContext will return a fired object or an error if context is done
func (w *Waiter) WaitWithContext(ctx context.Context) interface{} {
	w.waitCalled.Set(true)
	defer w.waitCalled.Set(false)

	var ret interface{}
	select {
	case <-ctx.Done():
		ret = errors.New("")
	case obj := <-w.c:
		ret = obj
	}

	return ret
}

// Fire will fire obj through the wait function if that function was called (and waiting), if not it will not fire.
func (w *Waiter) Fire(obj interface{}) {
	if w.waitCalled.Get() {
		w.c <- obj
	} else {
		w.queue.Add(obj) // will add if there is capacity
	}
}
