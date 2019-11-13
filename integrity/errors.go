package integrity

import "errors"

var (
	// ID
	ErrInvalidID = errors.New("the ID is NOT AN ID TYPE")

	// CRUD
	errInvalidCRUD = errors.New("the TYPE of operation is NOT ANY CRUD")
)
