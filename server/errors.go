package server

import "errors"

var (
	errNoTable  = errors.New("TABLE is NOT GIVEN in the request body")
	errNoBranch = errors.New("BRANCH is NOT GIVEN in the request body")
	errNoColumn = errors.New("COLUMN is NOT GIVEN in the request body")
)
