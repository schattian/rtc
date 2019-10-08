package git

type PullRequest struct {
	ID      int
	Name    string
	Commits []*Commit
}
