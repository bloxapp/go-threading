package queue

import (
	"sync"
)

// FuncQueue is the interface for managing a queue of Items
type FuncQueue interface {
	Add(item *Item) bool
	Pop() *Item
	ClearAndStop()
}

// Item represent some function
type Item struct {
	F func()
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
		queue: New(capacity),
		lock:  sync.Mutex{},
	}
	return &q
}

// Add will add an an item to the funcQueue, thread safe.
func (q *funcQueue) Add(e *Item) bool {
	return q.queue.Add(e)
}

// Pop will return and delete an an item from the funcQueue, thread safe.
func (q *funcQueue) Pop() *Item {
	ret, _ := q.queue.Pop().(*Item)
	return ret
}

// ClearAndStop will clear the funcQueue disable adding more items to it, thread safe.
func (q *funcQueue) ClearAndStop() {
	for _, item := range q.queue.All() {
		item.(*Item).Cancel()
	}
	q.queue.ClearAndStop()
}
