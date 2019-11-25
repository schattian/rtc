package integrity

import "errors"

var (
	// CRUD
	errInvalidCRUD = errors.New("the TYPE of operation is NOT ANY CRUD")
)
