package git

// setId wraps the id assignation to implement Storable interface
func (comm *Commit) setId(id int64) {
	comm.Id = id
}

// table returns the sql table name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (comm *Commit) table() string {
	return "commits"
}

// columns returns the columns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (comm *Commit) columns() []string {
	return []string{
		"id",
		"errored",
		"change_ids",
	}
}
