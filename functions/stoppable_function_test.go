package functions

import (
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/require"
)

func TestStoppableFuncLeaks(t *testing.T) {
	fStoppable := func(stopper FuncManager) (interface{}, error, bool) {
		for {
			if stopper.IsStopped() {
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		return nil, nil, false
	}

	fCompleted := func(stopper FuncManager) (interface{}, error, bool) {
		time.Sleep(time.Millisecond * 10)
		return nil, nil, false
	}

	t.Run("start with stop", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			sf := NewStoppableF(fStoppable)
			go func(sf *StoppableFunc) {
				time.Sleep(time.Millisecond * 10)
				sf.Manager.Stop()
			}(sf)
			sf.Start()
			wg.Done()
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})

	t.Run("start with completed", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			sf := NewStoppableF(fCompleted)
			sf.Start()
			wg.Done()
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})
}

func TestNewStoppableF(t *testing.T) {
	f := func(stopper FuncManager) (interface{}, error, bool) {
		time.Sleep(time.Millisecond * 50)
		return "done", nil, true
	}

	sf := NewStoppableF(f)
	res := sf.Start()
	require.EqualValues(t, "done", res.Obj.(string))
	require.NoError(t, res.Err)
	require.True(t, res.Completed)
}

func TestStoppedF(t *testing.T) {
	f := func(stopper FuncManager) (interface{}, error, bool) {
		for {
			if stopper.IsStopped() {
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		return nil, nil, false
	}

	sf := NewStoppableF(f)

	go func() {
		time.Sleep(time.Millisecond * 30)
		sf.Manager.Stop()
	}()

	res := sf.Start()
	require.Nil(t, res.Obj)
	require.NoError(t, res.Err)
	require.False(t, res.Completed)
}
