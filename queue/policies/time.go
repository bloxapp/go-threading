package policies

import "time"

var TimeOutPolicy = func(i ...interface{}) func() Policy {
	d := i[0].(time.Duration)
	return func() Policy {
		return NewTimePolicy(d)
	}
}

// TimePolicy evacuates queue items after a duration
type TimePolicy struct {
	t time.Time
}

func NewTimePolicy(d time.Duration) Policy {
	return &TimePolicy{
		t: time.Now().Add(d),
	}
}

func (tp *TimePolicy) Evacuate() bool {
	return time.Now().After(tp.t)
}
