package schema

import "errors"

var (
	// Table errs
	errNonexistentTable = errors.New("the table given doesn't exists")
	errForeignTable     = errors.New("the table given doesn't belongs to the given schema")

	// Column errs
	errNonexistentColumn = errors.New("the column given doesn't exists")
	errForeignColumn     = errors.New("the column given doesn't belongs to the given table")
)
