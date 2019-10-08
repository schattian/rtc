package git

type Owner struct {
	Collaborators []Collaborator
}

// Orchestrate sends the order to all the collaborators to execute
// the needed actions in order to achieve the commitment
func (o *Owner) Orchestrate(comm *Commit, name string, strategy changesMatcher) {
	pR := PullRequest{Name: name}
	for _, changes := range comm.GroupBy(strategy) {
		comm := &Commit{Changes: changes}
		pR.Commits = append(pR.Commits, comm)
	}
	o.Merge(pR)
	return
}

func (o *Owner) Merge(pR *PullRequest) {
	for _, comm := range pR.Commits {
	}
}

// ! A PR will be generated after a commit can be splitted into non-compatible changes
// ! A PR will be splitted in subCommits

// func (o *Owner) Pull(context.Context, *Commit) (*Commit, error) {
// 	return nil
// }

// func (o *Owner) Push(context.Context, *Commit) (Summary, error) {
// 	return
// }

// func (o *Owner) Delete(context.Context, *Commit) (Summary, error) {
// 	return
// }
