package git

type Result struct {
	CommitID int   `json:"commit_id,omitempty"`
	Error    error `json:"error,omitempty"`
}
