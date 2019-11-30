package git

// SetId wraps the id assignation to implement Storable interface
func (chg *Change) SetId(id int64) {
	chg.Id = id
}

// SQLTable returns the sql SQLTable name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (chg *Change) SQLTable() string {
	return "changes"
}

// SQLColumns returns the SQLColumns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (chg *Change) SQLColumns() []string {
	return []string{
		"id",
		"table_name",
		"column_name",
		"value_type",
		"str_value",
		"int_value",
		"float_32_value",
		"float_64_value",
		"json_value",
		"bytes_value",
		"entity_id",
		"index_id",
		"type",
		"options",
		"commited",
	}
}
