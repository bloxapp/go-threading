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
	// Len will return the number of items in the queue
	Len() int
	// ClearAndStop will clear the queue and stop it from adding more items
	ClearAndStop()
}

type Direction string

const (
	FIFO Direction = "FIFO"
	LIFO Direction = "LIFO"
)

// queue thread safe implementation of Queue
type queue struct {
	stop      bool
	queue     []Policy
	policies  []ApplyPolicy
	lock      sync.RWMutex
	capacity  int
	direction Direction
}

// New returns a new instance of funcQueue
func New(direction Direction, capacity int, policies ...ApplyPolicy) Queue {
	q := queue{
		queue:     make([]Policy, 0),
		lock:      sync.RWMutex{},
		capacity:  capacity,
		direction: direction,
		policies:  policies,
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

	policies := make([]Policy, 0)
	for _, p := range q.policies {
		policies = append(policies, p())
	}

	q.queue = append(q.queue, NewPolicyManager(e, policies))
	return true
}

// Pop will return and delete an an item from the funcQueue, thread safe.
func (q *queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return nil
	}

	var retItem interface{}
loop:
	for {
		qLen := len(q.queue)
		if qLen == 0 {
			break
		}

		if q.direction == FIFO {
			ret := q.queue[0]
			q.queue = q.queue[1:qLen]
			if !ret.Evacuate() {
				retItem = ret.Item()
				break loop
			}
		} else { // LIFO
			ret := q.queue[qLen-1]
			q.queue = q.queue[0 : qLen-1]
			if !ret.Evacuate() {
				retItem = ret.Item()
				break loop
			}
		}
	}

	return retItem
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

func (q *queue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.queue)
}

// ClearAndStop will clear the funcQueue disable adding more items to it, thread safe.
func (q *queue) ClearAndStop() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.stop = true
	q.queue = make([]Policy, 0)
}
