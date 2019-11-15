package git

// Result is a commitment result
type Result struct {
	CommitID int64 `json:"commit_id,omitempty"`
	Error    error `json:"error,omitempty"`
}
