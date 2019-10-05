package cmd

import "fmt"

// ProcessSignal is the wrapper of the synchronization process given a signal.
func ProcessSignal(signal Signal) (err error) {
	actions, err := signal.Decompose()
	if err != nil {
		return err
	}
	// errCh := make(chan error, len(actions))
	for _, action := range actions {
		// The sync process is did synchronously to avoid race conditions due to the unchained nature of the subscriptions obj
		if action {
			// Sync(errCh)
		} else {
			// Unsync(errCh)
		}
	}
	return nil
}

// The Signal is the representation of any possible action desired to pass to the wrapper.
// Those values can be: POST, PUT, DELETE, following the HTTP Verbs passed on the REST API communication.
// Notice that the act will handled by the wrapper, and the type itself doesn't complain about it's usability.
type Signal string

// Decompose will gives us a slice of boolean representation of the signal
// (useful to split complex signals into binary actions).
// Notice that the boolean of the returned slice of booleans matters, so negatives will be always first.
// false: a negative signal
// true: iota
func (s Signal) Decompose() ([]bool, error) {
	switch string(s) {
	case "POST":
		return []bool{true}, nil
	case "DELETE":
		return []bool{false}, nil
	case "PUT":
		return []bool{false, true}, nil
	default:
		return nil, fmt.Errorf("the signal \"%v\" cannot be decomoposed", s)
	}
}
