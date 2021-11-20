package queue

import "time"

type ApplyPolicy func() Policy

var TimeOutPolicy = func(i ...interface{}) func() Policy {
	d := i[0].(time.Duration)
	return func() Policy {
		return NewTimePolicy(d)
	}
}

type Policy interface {
	// Evacuate returns true if a msg should be evacuated from a queue
	Evacuate() bool
	Item() interface{}
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

func (tp *TimePolicy) Item() interface{} {
	return nil
}

// PolicyManager holds several policies and an item
type PolicyManager struct {
	policies []Policy
	item     interface{}
}

func NewPolicyManager(item interface{}, policies []Policy) Policy {
	return &PolicyManager{
		item:     item,
		policies: policies,
	}
}

func (m *PolicyManager) Evacuate() bool {
	for _, p := range m.policies {
		if p.Evacuate() {
			return true
		}
	}
	return false
}

func (m *PolicyManager) Item() interface{} {
	return m.item
}
