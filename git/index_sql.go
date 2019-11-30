package git

// SetId wraps the id assignation to implement Storable interface
func (idx *Index) SetId(id int64) {
	idx.Id = id
}

// SQLTable returns the sql SQLTable name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (idx *Index) SQLTable() string {
	return "indices"
}

// SQLColumns returns the SQLColumns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (idx *Index) SQLColumns() []string {
	return []string{
		"id",
	}
}
