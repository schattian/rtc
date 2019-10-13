package schema

import "errors"

var (
	// Table errs
	errNonexistentTable = errors.New("the TABLE given does NOT EXISTS")
	errForeignTable     = errors.New("the TABLE given does NOT BELONGS to the given SCHEMA")
	errInvalidOptionKey = errors.New("the provided OPTION KEY is INVALID OR NOT DEFINED")

	// Column errs
	errNonexistentColumn = errors.New("the COLUMN given does NOT EXISTS")
	errForeignColumn     = errors.New("the COLUMNS given does NOT BELONGS to the given TABLE")
)
