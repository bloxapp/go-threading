package timer

import (
	"go-threading/channel"
	"sync"
	"time"
)

const (
	waiterTimout = time.Millisecond * 50
)

// RoundTimer is a wrapper around timer to fit the use in an iBFT instance
type RoundTimer struct {
	timer     *time.Timer
	internalC *channel.Channel // internalC triggers when lapsed or cancelled
	resC      *channel.Channel

	stopped  bool
	syncLock sync.RWMutex
}

// New returns a new instance of RoundTimer
func New() *RoundTimer {
	return &RoundTimer{
		timer:     nil,
		internalC: channel.New(),
		resC:      channel.New(),
		stopped:   true,
		syncLock:  sync.RWMutex{},
	}
}

// ResultChan returns the result chan
// true if the timer lapsed or false if it was stopped
func (t *RoundTimer) ResultChan() *channel.Waiter {
	t.syncLock.Lock()
	defer t.syncLock.Unlock()
	return t.resC.Register()
}

// Reset will return a channel that sends true if the timer lapsed or false if it was cancelled
// If Start is called more than once, the first timer and chan are returned and used
func (t *RoundTimer) Reset(d time.Duration) {
	t.syncLock.Lock()
	defer t.syncLock.Unlock()

	t.stopped = false

	if t.timer != nil {
		// timer is already running, reset it.
		t.timer.Stop()
		t.timer.Reset(d)
	} else {
		// no running timer, create a new one
		go t.eventLoop()
		t.timer = time.AfterFunc(d, func() {
			t.syncLock.Lock()
			defer t.syncLock.Unlock()
			t.internalC.FireToAll(true)
			t.timer = nil
			t.stopEventLoop()
		})
	}
}

// Stopped returns true if there is no running timer
func (t *RoundTimer) Stopped() bool {
	t.syncLock.RLock()
	defer t.syncLock.RUnlock()
	return t.stopped
}

// Kill will stop the timer (without the ability to restart it) and send false on the result chan
func (t *RoundTimer) Kill() {
	t.syncLock.Lock()

	if t.timer != nil {
		t.timer.Stop()
	}
	t.stopEventLoop()

	t.syncLock.Unlock()

	t.fireChannelEvent(false)
}

func (t *RoundTimer) fireChannelEvent(value bool) {
	t.syncLock.RLock()
	defer t.syncLock.RUnlock()

	t.resC.FireToAll(value)
}

func (t *RoundTimer) stopEventLoop() {
	t.stopped = true
	t.internalC.FireToAll(false)
}

func (t *RoundTimer) eventLoop() {
	waiter := t.internalC.Register()
loop:
	for {
		res := waiter.Wait()
		if res.(bool) { // lapsed
			t.fireChannelEvent(true)
		} else { // killed
			break loop
		}
	}
}
