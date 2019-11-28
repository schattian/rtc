package git

import (
	"github.com/sebach1/rtc/integrity"
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

func (pR *PullRequest) copy() *PullRequest {
	if pR == nil {
		return nil
	}
	newPr := new(PullRequest)
	*newPr = *pR
	var newComms []*Commit
	for _, comm := range pR.Commits {
		newComms = append(newComms, comm.copy())
	}
	newPr.Commits = newComms
	newPr.Team = pR.Team.copy()
	return newPr
}
