package queue

import (
	"sync"
)

// FuncQueue is the interface for managing a queue of Items
type FuncQueue interface {
	Add(item *Item, indexes ...Index) bool
	Pop(index Index) *Item
}

// Item represent some function
type Item struct {
	F      func()
	Cancel func()
}

// funcQueue thread safe implementation of Queue
type funcQueue struct {
	stop  bool
	queue Queue
	lock  sync.Mutex
}

// NewFuncQueue returns a new instance of funcQueue
func NewFuncQueue(capacity int) FuncQueue {
	q := funcQueue{
		queue: New(FIFO, capacity),
		lock:  sync.Mutex{},
	}
	return &q
}

// Add will add an an item to the funcQueue, thread safe.
func (q *funcQueue) Add(e *Item, indexes ...Index) bool {
	return q.queue.Add(e, indexes...)
}

// Pop will return and delete an an item from the funcQueue, thread safe.
func (q *funcQueue) Pop(index Index) *Item {
	ret, _ := q.queue.Pop(index).(*Item)
	return ret
}
