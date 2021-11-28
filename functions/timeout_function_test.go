package functions

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestTimeoutFuncLeaks(t *testing.T) {
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
		for i := 0; i < 1; i++ {
			wg.Add(1)
			NewTimeoutF(fStoppable, time.Millisecond*10)
			wg.Done()
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})

	t.Run("start with completed", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			NewTimeoutF(fCompleted, time.Millisecond*1000)
			wg.Done()
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})
}

func TestCompletedInTime(t *testing.T) {
	f := func(stopper FuncManager) (interface{}, error, bool) {
		completed := false
		for i := 0; i <= 10; i++ {
			if stopper.IsStopped() {
				break
			}
			if i == 10 {
				completed = true
			}
			time.Sleep(time.Millisecond * 10)
		}
		if !completed {
			return nil, nil, completed
		}
		return "done", nil, completed
	}

	res := NewTimeoutF(f, time.Millisecond*1000)
	require.EqualValues(t, "done", res.Obj.(string))
	require.NoError(t, res.Err)
	require.True(t, res.Completed)
}

func TestDidNotCompleteInTime(t *testing.T) {
	f := func(stopper FuncManager) (interface{}, error, bool) {
		completed := false
		for i := 0; i <= 100; i++ {
			if stopper.IsStopped() {
				break
			}
			if i == 100 {
				completed = true
			}
			time.Sleep(time.Millisecond * 10)
		}
		return nil, nil, completed
	}

	res := NewTimeoutF(f, time.Millisecond*300)
	require.Nil(t, res.Obj)
	require.NoError(t, res.Err)
	require.False(t, res.Completed)
}
