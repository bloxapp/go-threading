package queue

import "sync"

// Queue is the interface for managing a queue of items
type Queue interface {
	Add(interface{}) bool
	Pop() interface{}
	All() []interface{}
	Len() int
	ClearAndStop()
}

// queue thread safe implementation of Queue
type queue struct {
	stop  bool
	queue []interface{}
	lock  sync.RWMutex
	capacity int
}

// New returns a new instance of funcQueue
func New(capacity int) Queue {
	q := queue{
		queue: make([]interface{}, 0),
		lock:  sync.RWMutex{},
		capacity: capacity,
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
		ret := q.queue[0]
		q.queue = q.queue[1:len(q.queue)]
		return ret
	}
	return nil
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

