package git

// GetId wraps the id retrieval to implement Storable interface
func (pR *PullRequest) GetId() int64 {
	return pR.Id
}

// SetId wraps the id assignation to implement Storable interface
func (pR *PullRequest) SetId(id int64) {
	pR.Id = id
}

// SQLTable returns the sql SQLTable name of the entity
//
// Testing: tested by using naming conventions. See internal/name pkg
func (pR *PullRequest) SQLTable() string {
	return "pull_requests"
}

// SQLColumns returns the SQLColumns each field represent on db
// Notice the returned slice is the list of struct tags of exported fields
// It's done to avoid reflection
//
// Testing: tested by using reflection at Columns_Test to check being the tags
func (pR *PullRequest) SQLColumns() []string {
	return []string{
		"id",
	}
}
