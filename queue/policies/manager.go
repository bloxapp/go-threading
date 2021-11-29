package policies

type PolicyManager interface {
	Evacuate() bool
	AddPolicy(policy Policy)
}

// policyManager holds several policies and an item
type policyManager struct {
	policies []Policy
}

// NewPolicyManager returns a policy manager instance
func NewPolicyManager(policies []Policy) PolicyManager {
	return &policyManager{
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

func (m *policyManager) AddPolicy(policy Policy) {
	m.policies = append(m.policies, policy)
}
