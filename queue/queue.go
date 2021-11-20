package queue

import (
	"go-threading/channel"
	policies2 "go-threading/queue/policies"
	"sync"
	"time"
)

const (
	PopWaitSleepTime = time.Millisecond * 50
)

// Queue is the interface for managing a queue of items
// An item can have several policies which dictate when the item is evicted from the queue.
// Items are evicted (if need be) when the queue reaches capacity and a new item needs to be added
// When adding an item an array of indexes can be provided for the item, a pop call needs to be called for each index later on.
type Queue interface {
	// Add will add an item to the queue. If one or more indexes are provided will store accordingly, otherwise the default index will be used
	Add(interface{}, ...Index) bool
	// Pop will return the next item or nil. If no index provided, the default index will be used
	Pop(Index) interface{}
	// PopWait returns a waiter which can be used to pop or wait for a new object and then pop. If no index provided, the default index will be used.
	PopWait(Index) *channel.Waiter
	// Len will return the number of items in the queue
	Len() int
}

type Direction string

const (
	FIFO Direction = "FIFO"
	LIFO Direction = "LIFO"
)

type Index string

const (
	DefaultItemIndex Index = "DefaultItemIndex"
)

// queue thread safe implementation of Queue
type queue struct {
	stop      bool
	queue     map[Index][]policies2.PolicyManager
	policies  []policies2.ApplyPolicy
	lock      sync.RWMutex
	capacity  int
	direction Direction
	count     int
}

// New returns a new instance of funcQueue
func New(direction Direction, capacity int, policies ...policies2.ApplyPolicy) Queue {
	q := queue{
		queue:     make(map[Index][]policies2.PolicyManager),
		lock:      sync.RWMutex{},
		capacity:  capacity,
		direction: direction,
		policies:  policies,
		count:     0,
	}
	return &q
}

// Add will add an item to the queue, thread safe.
func (q *queue) Add(e interface{}, indexes ...Index) bool {
	if !q.preAddCheck(indexes...) {
		return false
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return false
	}

	policies := make([]policies2.Policy, 0)
	for _, p := range q.policies {
		policies = append(policies, p())
	}

	if len(indexes) > 0 {
		for _, index := range indexes {
			if q.queue[index] == nil {
				q.queue[index] = make([]policies2.PolicyManager, 0)
			}
			q.queue[index] = append(q.queue[index], policies2.NewPolicyManager(e, policies))
			q.count++
		}
	} else {
		if q.queue[DefaultItemIndex] == nil {
			q.queue[DefaultItemIndex] = make([]policies2.PolicyManager, 0)
		}
		q.queue[DefaultItemIndex] = append(q.queue[DefaultItemIndex], policies2.NewPolicyManager(e, policies))
		q.count++
	}

	return true
}

// Pop will return and delete an item from the funcQueue, thread safe.
func (q *queue) Pop(index Index) interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.stop {
		return nil
	}

	if q.evictItems() == 0 || len(q.queue[index]) == 0 {
		return nil
	}

	indexedQ := q.queue[index]
	qLen := len(indexedQ)
	var ret policies2.PolicyManager
	if q.direction == FIFO {
		ret = indexedQ[0]
		q.queue[index] = indexedQ[1:qLen]
	} else { // LIFO
		ret = indexedQ[qLen-1]
		q.queue[index] = indexedQ[0 : qLen-1]
	}

	// delete index if empty
	if len(q.queue[index]) == 0 {
		delete(q.queue, index)
	}

	return ret.Item()
}

func (q *queue) PopWait(index Index) *channel.Waiter {
	c := channel.New()
	go func() {
	loop:
		for {
			if obj := q.Pop(index); obj != nil {
				c.FireToAll(obj)
				break loop
			}
			time.Sleep(PopWaitSleepTime)
		}
	}()
	return c.Register()
}

func (q *queue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.count
}

// evictItems evicts items according to policy and returns total (after eviction) count
// not thread safe, should be called safely
func (q *queue) evictItems() int {
	newCount := 0
	for index, indexedQ := range q.queue {
		newQ := make([]policies2.PolicyManager, 0)
		for _, i := range indexedQ {
			if !i.Evacuate() {
				newQ = append(newQ, i)
				newCount++
			}
		}
		q.queue[index] = newQ
	}
	q.count = newCount
	return newCount
}

// preAddCheck will return true if possible to add item
func (q *queue) preAddCheck(indexes ...Index) bool {
	itemsToAdd := 1
	if len(indexes) > 0 {
		itemsToAdd = len(indexes)
	}

	if q.Len()+itemsToAdd > q.capacity {
		q.lock.Lock()
		l := q.evictItems()
		q.lock.Unlock()
		if l+itemsToAdd > q.capacity {
			return false
		}
	}
	return true
}
