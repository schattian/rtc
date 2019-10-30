package integrity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrorsSeparator is the expected string to use when stringifiying multiple errors to one
const ErrorsSeparator = "; "

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

var vErrRegex = regexp.MustCompile(`validation error: (.*). At`)

// UnwrapValidationError returns a slice of the errs given by a validation error
// - Verbose: arg specifies if want to get details about where the error occured
func UnwrapValidationError(err error, verbose bool) (errs []error) {
	if err == nil {
		return
	}
	strErr := err.Error()
	strErrs := strings.Split(strErr, ErrorsSeparator)
	strErrs = strErrs[0 : len(strErrs)-1] // Trims the latest separator
	for _, strErr := range strErrs {
		if !verbose {
			strErr = vErrRegex.FindAllString(strErr, 1)[0]
			strErr = strings.TrimSuffix(strings.TrimPrefix(strErr, "validation error: "), ". At")
		}
		errs = append(errs, errors.New(strErr))
	}
	return
}
