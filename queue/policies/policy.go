package policies

type ApplyPolicy func() Policy

type Policy interface {
	// Evacuate returns true if a msg should be evacuated from a queue
	Evacuate() bool
}
