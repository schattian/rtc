package git

// SetId wraps the id assignation to implement Storable interface
func (comm *Commit) SetId(id int64) {
	comm.Id = id
}

// Table returns the sql table name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (comm *Commit) Table() string {
	return "commits"
}

// Columns returns the columns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (comm *Commit) Columns() []string {
	return []string{
		"id",
		"errored",
		"change_ids",
	}
}
