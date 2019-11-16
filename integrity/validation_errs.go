package integrity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidationError is the custom validation error which stores:
// - OriginType: the struct which gave the error
// - Err: the err to be wrapped by .Error()
type ValidationError struct {
	OriginType string
	OriginName string
	Err        error
}

func (vErr *ValidationError) Error() string {
	base := vErr.OriginType + " validation error: " + vErr.Err.Error()
	if vErr.OriginName != "" {
		base += fmt.Sprintf(". At %v", vErr.OriginName)
	}
	return strings.TrimSpace(base)
}

var vErrRegex = regexp.MustCompile(`validation error: (.*). At`)

// UnwrapValidationError returns a slice of the errs given by a validation error
// - Verbose: arg specifies if want to get details about where the error occured (vErr.Origin)
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
