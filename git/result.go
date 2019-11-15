package git

// Result is a commitment result
type Result struct {
	CommitId int64 `json:"commit_id,omitempty"`
	Error    error `json:"error,omitempty"`
}
