package channel

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const QueueSize = 5

var ContextDoneErr = errors.New("WAITER_CONTEXT_DONE")

type Waiter struct {
	id string
	c  chan interface{}
}

func NewWaiter() *Waiter {
	return &Waiter{
		id: uuid.New().String(),
		c:  make(chan interface{}, QueueSize),
	}
}

// Wait will block until a new obj is passed from a queue or from firing.
// If queue has items, will return immediately after popping firs item
func (w *Waiter) Wait() interface{} {
	return <-w.c
}

// WaitWithTimeout will return a fired object or an error if deadline exceeded
func (w *Waiter) WaitWithTimeout(duration time.Duration) interface{} {
	c, _ := context.WithTimeout(context.Background(), duration)
	return w.WaitWithContext(c)
}

// WaitWithContext will return a fired object or an error if context is done
func (w *Waiter) WaitWithContext(ctx context.Context) interface{} {
	var ret interface{}
	select {
	case <-ctx.Done():
		ret = ContextDoneErr
	case obj := <-w.c:
		ret = obj
	}

	return ret
}

// Fire will fire obj through the wait function if that function was called (and waiting), if not it will not fire.
func (w *Waiter) Fire(obj interface{}) {
	w.c <- obj
}
