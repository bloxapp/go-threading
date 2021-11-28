package timer

import (
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/require"
)

func TestLeaks(t *testing.T) {
	t.Run("elapsed", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			timer := New()
			go func(timer *RoundTimer) {
				timer.ResultChan().Wait()
				wg.Done()
			}(timer)
			timer.Reset(time.Millisecond * 25)
		}
		wg.Wait()
		time.Sleep(time.Millisecond * 500) // wait for all go routines to close
		goleak.VerifyNone(t)
	})

	t.Run("reset twice and elapse", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			timer := New()
			go func(timer *RoundTimer) {
				timer.ResultChan().Wait()
				wg.Done()
			}(timer)
			timer.Reset(time.Millisecond * 150)
			time.AfterFunc(time.Millisecond*10, func() {
				timer.Reset(time.Millisecond * 15)
			})
		}
		wg.Wait()
		time.Sleep(time.Millisecond * 500) // wait for all go routines to close
		goleak.VerifyNone(t)
	})

	t.Run("killed", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			timer := New()
			go func(timer *RoundTimer) {
				time.Sleep(time.Millisecond * 25)
				timer.Kill()
				wg.Done()
			}(timer)
		}
		wg.Wait()
		time.Sleep(time.Millisecond * 500) // wait for all go routines to close
		goleak.VerifyNone(t)
	})
}

func TestRoundTimer_Reset(t *testing.T) {
	timer := New()
	timer.Reset(time.Millisecond * 100)
	require.False(t, timer.Stopped())
	res := timer.ResultChan().Wait()
	require.True(t, res.(bool))
	require.True(t, timer.Stopped())

	timer.Reset(time.Millisecond * 100)
	require.False(t, timer.Stopped())
	res = timer.ResultChan().Wait()
	require.True(t, res.(bool))
	require.True(t, timer.Stopped())
}

func TestRoundTimer_ResetTwice(t *testing.T) {
	timer := New()
	timer.Reset(time.Millisecond * 100)
	timer.ResultChan().Wait()
	timer.Reset(time.Millisecond * 100)
	res := timer.ResultChan().Wait()
	require.True(t, res.(bool))
}

func TestRoundTimer_ResetBeforeLapsed(t *testing.T) {
	timer := New()
	timer.Reset(time.Millisecond * 100)
	timer.Reset(time.Millisecond * 300)

	t1 := time.Now()
	res := timer.ResultChan().Wait()
	t2 := time.Since(t1)
	require.True(t, res.(bool))
	require.Greater(t, t2.Milliseconds(), (time.Millisecond * 150).Milliseconds())
}

func TestRoundTimer_Stop(t *testing.T) {
	timer := New()
	timer.Reset(time.Millisecond * 500)
	go func() {
		time.Sleep(time.Millisecond * 100)
		timer.Kill()
	}()
	require.False(t, timer.ResultChan().Wait().(bool))
}
