package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bloxapp/go-threading/threadsafe"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/require"
)

func TestChannelLeaks(t *testing.T) {
	t.Run("wait and fire", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			c := New()
			go func(c *Channel) {
				c.FireToAll("test")
			}(c)
			c.Register().Wait()
			wg.Done()
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})

	t.Run("wait with timeout, don't fire", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			c := New()
			c.Register().WaitWithTimeout(time.Millisecond * 10)
			wg.Done()
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})
}

func TestChannel_RegisterAndFire(t *testing.T) {
	c := New()

	fired := make([]*threadsafe.SafeBool, 0)
	for i := 0; i < 100; i++ {
		w := c.Register()
		firedBool := threadsafe.Bool()
		fired = append(fired, firedBool)
		go func(w *Waiter, firedBool *threadsafe.SafeBool) {
			w.Wait()
			firedBool.Set(true)
		}(w, firedBool)
	}

	time.Sleep(time.Millisecond * 25)
	c.FireToAll(true)

	// verify
	for i, b := range fired {
		t.Run(fmt.Sprintf("waiter: %d", i), func(t *testing.T) {
			require.True(t, b.Get())
		})
	}
}

func TestChannel_DeRegister(t *testing.T) {
	c := New()

	waiters := make([]*Waiter, 0)
	for i := 0; i < 100; i++ {
		waiters = append(waiters, c.Register())
	}

	// deregister and check
	for _, w := range waiters {
		c.DeRegister(w)
	}
	require.Len(t, c.registers, 0)
}

func TestChannel_FireOnceToAll(t *testing.T) {
	c := New()

	fired := make([]*threadsafe.AnyObj, 0)
	for i := 0; i < 100; i++ {
		w := c.Register()
		firedObj := threadsafe.Any()
		fired = append(fired, firedObj)
		go func(w *Waiter, firedObj *threadsafe.AnyObj) {
			obj := w.Wait()
			require.True(t, obj.(bool))
			obj = w.Wait()
			firedObj.Set(obj)
		}(w, firedObj)
	}
	time.Sleep(time.Millisecond * 100)
	c.FireOnceToAll(true)
	time.Sleep(time.Millisecond * 100)

	// verify
	for i, b := range fired {
		t.Run(fmt.Sprintf("waiter: %d", i), func(t *testing.T) {
			require.EqualValues(t, b.Get(), ChannelClosed)
		})
	}
}
