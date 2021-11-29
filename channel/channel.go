package channel

import (
	"sync"

	"github.com/bloxapp/go-threading/threadsafe"
)

const ChannelClosed = "channel_closed"

type Channel struct {
	lock      sync.RWMutex
	registers map[string]*Waiter
	cancelled *threadsafe.SafeBool
}

func New() *Channel {
	return &Channel{
		lock:      sync.RWMutex{},
		registers: make(map[string]*Waiter),
		cancelled: threadsafe.Bool(),
	}
}

// Register will return a waiter if channel is active, nil if it's cancelled
func (c *Channel) Register() *Waiter {
	c.lock.Lock()
	defer c.lock.Unlock()

	ret := NewWaiter()
	c.registers[ret.id] = ret
	return ret
}

func (c *Channel) DeRegister(waiter *Waiter) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.registers[waiter.id]; ok {
		delete(c.registers, waiter.id)
	}
}

func (c *Channel) DeRegisterAll() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.registers = make(map[string]*Waiter)
}

// FireToAll will fire the object thorough the waiters if not cancelled
func (c *Channel) FireToAll(obj interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, w := range c.registers {
		w.Fire(obj)
	}
}

// FireOnceToAll will fire the object through the waiters if not cancelled, will cancel channel after
func (c *Channel) FireOnceToAll(obj interface{}) {
	c.FireToAll(obj)
	c.CancelAll()
}

// CancelAll will fire nil to all waiters and will not fire any obj again
func (c *Channel) CancelAll() {
	c.cancelled.Set(true)
	c.FireToAll(ChannelClosed)
}
