package git

import "github.com/sebach1/git-crud/internal/integrity"

type PullRequest struct {
	ID      int
	Team    *Team
	Commits []*Commit
}

func (pR *PullRequest) AssignTeam(community *Community, schName integrity.SchemaName) error {
	team, err := community.LookFor(schName)
	if err != nil {
		return err
	}
	pR.Team = team
	return nil
}
