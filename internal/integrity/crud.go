package integrity

import "fmt"

// CRUD is any string which preserves any state of Create Retrieve Update or Delete
type CRUD string

// Validate asserts if the string is any kind of CRUD actions
func (crud CRUD) Validate() error {
	for _, action := range []CRUD{"create", "retrieve", "update", "delete"} {
		if crud == action {
			return nil
		}
	}
	return fmt.Errorf("INVALID CRUD operation: %v", crud)
}
