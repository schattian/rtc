package git

import (
	"github.com/sebach1/git-crud/integrity"
)

// A PullRequest connects a group of Commits with a team
type PullRequest struct {
	Id      int
	Team    *Team
	Commits []*Commit
}

// AssignTeam looks up for a team given a schemaName and a community
// Notice that it cleans up the current Team
func (pR *PullRequest) AssignTeam(community *Community, schName integrity.SchemaName) error {
	pR.Team = &Team{}
	team, err := community.LookFor(schName)
	if err != nil {
		return err
	}
	pR.Team = team
	return nil
}
