package functions

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewStoppableF(t *testing.T) {
	f :=  func(stopper FuncManager) (interface{}, error, bool) {
		time.Sleep(time.Millisecond*50)
		return "done",nil,true
	}

	sf := NewStoppableF(f)
	res := sf.Start()
	require.EqualValues(t, "done", res.Obj.(string))
	require.NoError(t, res.Err)
	require.True(t, res.Completed)
}

func TestStoppedF(t *testing.T) {
	f :=  func(stopper FuncManager) (interface{}, error, bool) {
		for {
			if stopper.IsStopped(){
				break
			}
			time.Sleep(time.Millisecond*10)
		}
		return nil, nil, false
	}

	sf := NewStoppableF(f)

	go func() {
		time.Sleep(time.Millisecond*30)
		sf.Manager.Stop()
	}()

	res := sf.Start()
	require.Nil(t,res.Obj)
	require.NoError(t, res.Err)
	require.False(t, res.Completed)
}
