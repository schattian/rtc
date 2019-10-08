package integrity

import "fmt"

type CRUD string

func (crud CRUD) Validate() error {
	for _, action := range []CRUD{"create", "retrieve", "update", "delete"} {
		if crud == action {
			return nil
		}
	}
	return fmt.Errorf("INVALID CRUD operation: %v", crud)
}
