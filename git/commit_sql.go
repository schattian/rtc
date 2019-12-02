package git

// GetId wraps the id retrieval to implement Storable interface
func (comm *Commit) GetId() int64 {
	return comm.Id
}

// SetId wraps the id assignation to implement Storable interface
func (comm *Commit) SetId(id int64) {
	comm.Id = id
}

// SQLTable returns the sql SQLTable name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (comm *Commit) SQLTable() string {
	return "commits"
}

// SQLColumns returns the SQLColumns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (comm *Commit) SQLColumns() []string {
	return []string{
		"id",
		"errored",
	}
}
