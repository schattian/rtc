package git

import "errors"

var (
	// Changes
	errDuplicatedChg = errors.New("the change is ALREADY COMMITED")
	errZeroTable     = errors.New("change's TABLE cannot be ZERO")
	errZeroColumn    = errors.New("change's COLUMN cannot be ZERO")
)
