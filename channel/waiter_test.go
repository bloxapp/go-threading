package channel

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bloxapp/go-threading/threadsafe"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/require"
)

func TestWaiterLeaks(t *testing.T) {
	t.Run("no waiting called", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			w := NewWaiter()
			go func(w *Waiter) {
				time.Sleep(time.Millisecond * 25)
				w.Fire("test")
				wg.Done()
			}(w)
		}
		wg.Wait()
		goleak.VerifyNone(t)
	})

	t.Run("wait called", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			w := NewWaiter()
			go func(w *Waiter) {
				time.Sleep(time.Millisecond * 25)
				w.Fire("test")
			}(w)
			w.Wait()
			wg.Done()
		}
		wg.Wait()
		goleak.VerifyNone(t)
	})

	t.Run("wait with timeout", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			w := NewWaiter()
			w.WaitWithTimeout(time.Millisecond * 10)
			wg.Done()
		}
		wg.Wait()
		goleak.VerifyNone(t)
	})
}

func TestWaiterNormal(t *testing.T) {
	w := NewWaiter()

	res := threadsafe.Bool()
	go func(res *threadsafe.SafeBool) {
		newVal := w.Wait()
		res.Set(newVal.(bool))
	}(res)

	time.Sleep(time.Millisecond * 10)
	w.Fire(true)
	time.Sleep(time.Millisecond * 10)
	require.True(t, res.Get())
}

func TestWaiterQueue(t *testing.T) {
	w := NewWaiter()

	for i := 0; i < QueueSize-1; i++ {
		w.Fire(true)
	}
	w.Fire(nil)

	cnt := 0
	for {
		cnt++
		if w.Wait() == nil {
			break
		}
	}
	require.EqualValues(t, QueueSize, cnt)
}

func TestWaiterWithTimeout(t *testing.T) {
	w := NewWaiter()

	res := threadsafe.Bool()
	fired := threadsafe.Bool()
	go func(res *threadsafe.SafeBool) {
		newVal := w.WaitWithTimeout(time.Millisecond * 50)
		if obj, isBool := newVal.(bool); isBool {
			res.Set(obj)
		}
		fired.Set(true)
	}(res)

	time.Sleep(time.Millisecond * 60)
	require.True(t, fired.Get())
	require.False(t, res.Get())
}

func TestWaiterWithContext(t *testing.T) {
	w := NewWaiter()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Millisecond * 30)
		cancel()
	}()

	res := threadsafe.Bool()
	fired := threadsafe.Bool()
	go func(res *threadsafe.SafeBool) {
		newVal := w.WaitWithContext(ctx)
		if obj, isBool := newVal.(bool); isBool {
			res.Set(obj)
		}
		fired.Set(true)
	}(res)

	time.Sleep(time.Millisecond * 40)
	require.True(t, fired.Get())
	require.False(t, res.Get())
}
