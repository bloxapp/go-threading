package policies

var CancelledPolicy = func(i ...interface{}) func() Policy {
	return func() Policy {
		return NewCancelledPolicy()
	}
}

// cancelledPolicy evacuates queue items immediately
type cancelledPolicy struct {
}

func NewCancelledPolicy() Policy {
	return &cancelledPolicy{}
}

func (tp *cancelledPolicy) Evacuate() bool {
	return true
}
