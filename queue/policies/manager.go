package policies

type PolicyManager interface {
	Evacuate() bool
	Item() interface{}
}

// policyManager holds several policies and an item
type policyManager struct {
	policies []Policy
	item     interface{}
}

func NewPolicyManager(item interface{}, policies []Policy) PolicyManager {
	return &policyManager{
		item:     item,
		policies: policies,
	}
}

func (m *policyManager) Evacuate() bool {
	for _, p := range m.policies {
		if p.Evacuate() {
			return true
		}
	}
	return false
}

func (m *policyManager) Item() interface{} {
	return m.item
}
