package queue

import (
	"sync"
)

// Item represent some function
type Item struct {
	F func()
	Cancel func()
}

// ItemQueue is the interface for managing a queue of functions
type ItemQueue interface {
	Add(*Item) bool
	Pop() *Item
	ClearAndStop()
}

// queue thread safe implementation of ItemQueue
type queue struct {
	stop  bool
	queue []*Item
	lock  sync.Mutex
}

// New returns a new instance of queue
func New() ItemQueue {
	q := queue{
		queue: make([]*Item, 0),
		lock:  sync.Mutex{},
	}
	return &q
}

// Add will add an an item to the queue, thread safe.
func (q *queue) Add(e *Item) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return false
	}

	q.queue = append(q.queue, e)
	return true
}

// Pop will return and delete an an item from the queue, thread safe.
func (q *queue) Pop() *Item {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return nil
	}

	if len(q.queue) > 0 {
		ret := q.queue[0]
		q.queue = q.queue[1:len(q.queue)]
		return ret
	}
	return nil
}

// ClearAndStop will clear the queue disable adding more items to it, thread safe.
func (q *queue) ClearAndStop() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.stop = true
	for _, item := range q.queue {
		item.Cancel()
	}
	q.queue = make([]*Item, 0)
}
