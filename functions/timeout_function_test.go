package functions

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCompletedInTime(t *testing.T) {
	f :=  func(stopper FuncManager) (interface{}, error, bool) {
		completed := false
		for i:= 0 ; i <= 10 ; i++ {
			if stopper.IsStopped(){
				break
			}
			if i == 10 {
				completed = true
			}
			time.Sleep(time.Millisecond*10)
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
	f :=  func(stopper FuncManager) (interface{}, error, bool) {
		completed := false
		for i:= 0 ; i <= 100 ; i++ {
			if stopper.IsStopped(){
				break
			}
			if i == 100 {
				completed = true
			}
			time.Sleep(time.Millisecond*10)
		}
		return nil, nil, completed
	}

	res := NewTimeoutF(f, time.Millisecond*300)
	require.Nil(t, res.Obj)
	require.NoError(t, res.Err)
	require.False(t, res.Completed)
}