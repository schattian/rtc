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

func (vErr ValidationError) Error() string {
	if vErr.Err == nil {
		return ""
	}
	base := vErr.OriginType + " validation error: " + vErr.Err.Error()
	if vErr.OriginName != "" {
		base += fmt.Sprintf(". At %v", vErr.OriginName)
	}
	return strings.TrimSpace(base)
}

var fullVErrRegex = regexp.MustCompile(`validation error: (.*). At`)
var tinyVErrRegex = regexp.MustCompile(`validation error: (.*)`)

// UnwrapValidationError will undo .Error() boilerplate (skipping origin information, and giving .Err)
func UnwrapValidationError(vErr error) error {
	if vErr == nil || vErr.Error() == "" {
		return nil
	}
	strErr := vErr.Error()
	var vErrRegex *regexp.Regexp
	if strings.Contains(strErr, ". At") {
		vErrRegex = fullVErrRegex
	} else {
		vErrRegex = tinyVErrRegex
	}
	strErr = vErrRegex.FindAllString(strErr, 1)[0]
	strErr = strings.TrimSuffix(strings.TrimPrefix(strErr, "validation error: "), ". At")
	return errors.New(strErr)
}
