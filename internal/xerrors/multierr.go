package xerrors

// ErrorsSeparator is the expected string to use when stringifiying multiple errors to one
const ErrorsSeparator = "; "

// MultiErr is an error interface with multiple error
// .Error() will retrieve the appended errs separated by ErrorsSeparator
type MultiErr []error

// UnwrapAll acts as a map func over the MultiErr entity, being the map func the proportioned unwrapper
func (errs MultiErr) UnwrapAll(unwrapper func(error) error) (unwrappedErrs []error) {
	for _, err := range errs {
		unwrappedErrs = append(unwrappedErrs, unwrapper(err))
	}
	return
}

// NewMultiErr returns a MultiErr with the given errs
func NewMultiErr(errs ...error) (mErr MultiErr) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		mErr = append(mErr, err)
	}
	return
}

func NewMultiErrFromCh(errCh chan error) (mErr MultiErr) {
	for err := range errCh {
		mErr = append(mErr, err)
	}
	return
}

func (errs MultiErr) Error() string {
	var strBaseErr string
	for _, err := range errs {
		if err == nil {
			continue
		}
		strBaseErr += err.Error()
		strBaseErr += ErrorsSeparator
	}
	return strBaseErr
}
