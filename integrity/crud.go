package integrity

// CRUD is any string which preserves any state of Create Retrieve Update or Delete
type CRUD string

// Validate asserts if the string is any kind of CRUD actions
func (crud CRUD) Validate() error {
	for _, action := range []CRUD{"create", "retrieve", "update", "delete"} {
		if crud == action {
			return nil
		}
	}
	return errInvalidCRUD
}

// ToHTTPVerb retrieves the CRUD synonym by REST HTTP convention
func (crud CRUD) ToHTTPVerb() string {
	switch crud {
	case "create":
		return "POST"
	case "retrieve":
		return "GET"
	case "update":
		return "PUT"
	case "delete":
		return "DELETE"
	}
	return ""
}
