package functions

import (
	"go-threading/channel"

	"github.com/pkg/errors"
)

// FuncWithStop is the interface of functions to trigger. Returns an optional object, error and bool (true if finished execution, false if stopped)
type FuncWithStop = func(stopper FuncManager) (interface{}, error, bool)

type FuncResult struct {
	Obj       interface{}
	Err       error
	Completed bool
}

type StoppableFunc struct {
	fn      FuncWithStop
	Manager FuncManager
	Result  *channel.Channel
}

// NewStoppableF will run a provided function in a new go routine with a funcManager and a results channel which returns FuncResult
func NewStoppableF(fn FuncWithStop) *StoppableFunc {
	return &StoppableFunc{
		fn:      fn,
		Manager: newFuncManager(),
		Result:  channel.New(),
	}
}

func (s *StoppableFunc) Start() *FuncResult {
	w := s.Result.Register()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.Result.FireToAll(&FuncResult{
					Err:       errors.Errorf("panic: %s", err),
					Completed: false,
				})
			}
		}()

		res, err, completed := s.fn(s.Manager)
		s.Result.FireToAll(&FuncResult{
			Obj:       res,
			Err:       err,
			Completed: completed,
		})
	}()

	res := w.Wait()
	return res.(*FuncResult)
}
