package git

// SetID wraps the id assignation to implement Storable interface
func (chg *Change) SetID(id int64) {
	chg.ID = id
}

// Table returns the sql table name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (chg *Change) Table() string {
	return "changes"
}

// Columns returns the columns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (chg *Change) Columns() []string {
	return []string{
		"id",
		"table_name",
		"column_name",
		"str_value",
		"int_value",
		"float32_value",
		"float64_value",
		"json_value",
		"bytes_value",
		"entity_id",
		"type",
		"options",
		"commited",
	}
}
