package schematypes

// Result is a changes PUSHing result, grouped by changesID that was sent grouped
// See Commit.GroupBy to see grouping details
type Result struct {
	ChangesID []int `json:"changes_id,omitempty"`
	Error     error `json:"error,omitempty"`
}

// Summary is the channel which stores async all the results
type Summary chan *Result
