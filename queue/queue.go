package queue

import (
	"go-threading/channel"
	"sync"
	"time"
)

const (
	PullSleepTime = time.Millisecond * 50
)

// Queue is the interface for managing a queue of items
type Queue interface {
	// Add will add an item to the queue
	Add(interface{}) bool
	// Pop will return the next item or nil
	Pop() interface{}
	// PopWait returns a waiter which can be used to pop or wait for a new object and then pop
	PopWait() *channel.Waiter
	// All returns all queue items
	All() []interface{}
	// Len will return the number of items in the queue
	Len() int
	// ClearAndStop will clear the queue and stop it from adding more items
	ClearAndStop()
}

type QueueDirection string

const (
	FIFO QueueDirection = "FIFO"
	LIFO QueueDirection = "LIFO"
)

// queue thread safe implementation of Queue
type queue struct {
	stop      bool
	queue     []interface{}
	lock      sync.RWMutex
	capacity  int
	direction QueueDirection
}

// New returns a new instance of funcQueue
func New(direction QueueDirection, capacity int) Queue {
	q := queue{
		queue:     make([]interface{}, 0),
		lock:      sync.RWMutex{},
		capacity:  capacity,
		direction: direction,
	}
	return &q
}

// Add will add an item to the queue, thread safe.
func (q *queue) Add(e interface{}) bool {
	if q.Len() >= q.capacity {
		return false
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return false
	}

	q.queue = append(q.queue, e)
	return true
}

// Pop will return and delete an an item from the funcQueue, thread safe.
func (q *queue) Pop() interface{} {
	qLen := q.Len() // called before lock as it's a locking call as well

	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return nil
	}

	if qLen > 0 {
		if q.direction == FIFO {
			ret := q.queue[0]
			q.queue = q.queue[1:qLen]
			return ret
		} else { // LIFO
			ret := q.queue[qLen-1]
			q.queue = q.queue[0 : qLen-1]
			return ret
		}

	}
	return nil
}

func (q *queue) PopWait() *channel.Waiter {
	c := channel.New()
	go func() {
	loop:
		for {
			if obj := q.Pop(); obj != nil {
				c.FireToAll(obj)
				break loop
			}
			time.Sleep(PullSleepTime)
		}
	}()
	return c.Register()
}

func (q *queue) All() []interface{} {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.queue
}

func (q *queue) Len() int {
	return len(q.All())
}

// ClearAndStop will clear the funcQueue disable adding more items to it, thread safe.
func (q *queue) ClearAndStop() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.stop = true
	q.queue = make([]interface{}, 0)
}
