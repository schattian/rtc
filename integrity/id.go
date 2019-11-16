package integrity

import "strings"

// Id is a string in order to handle with every kind of id types (UUID/GUID/int)
type Id string

// IsNil verifies if the id is zero-valued
func (id Id) IsNil() bool {
	withoutZeros := strings.ReplaceAll(string(id), "0", "")
	if withoutZeros == "" {
		return true
	}
	return false
}
