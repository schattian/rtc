package integrity

import "fmt"

// Validator is any func that returns an error if it doesn't pass the integrity checks for a given interface
type Validator func(interface{}) error

// ValidationError is the custom validation error which stores:
// - Origin: the struct which gave the error
// - Err: the err to be wrapped by .Error()
type ValidationError struct {
	Origin     string
	OriginName string
	Err        error
}

func (vErr *ValidationError) Error() string {
	return fmt.Sprintf("%v validation error: %v. At %v", vErr.Origin, vErr.Err.Error(), vErr.OriginName)

}
