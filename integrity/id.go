package integrity

// Id is a string in order to handle with every kind of id types (UUID/GUID/int)
type Id string

// IsNil verifies if the id is zero-valued
func (id Id) IsNil() bool {
	if id == "0" || id == "" {
		return true
	}
	return false
}
