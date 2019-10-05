package maps

import "errors"

var (
	// Table errs
	errUnexistantTable = errors.New("the table given doesnt exists")
	errForeignTable    = errors.New("the table given doesnt belongs to the given map")

	// Column errs
	errUnexistantColumn = errors.New("the column given doesnt exists")
	errForeignColumn    = errors.New("the column given doesnt belongs to the given table")
)
