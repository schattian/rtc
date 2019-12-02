package git

// GetId wraps the id retrieval to implement Storable interface
func (b *Branch) GetId() int64 {
	return b.Id
}

// SetId wraps the id assignation to implement Storable interface
func (b *Branch) SetId(id int64) {
	b.Id = id
}

// SQLTable returns the sql SQLTable name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (b *Branch) SQLTable() string {
	return "branches"
}

// SQLColumns returns the SQLColumns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (b *Branch) SQLColumns() []string {
	return []string{
		"id",
		"name",
		"index_id",
	}
}
