package git

import (
	"sync"

	"github.com/sebach1/git-crud/internal/integrity"
)

// A PullRequest connects a group of Commits with a team
type PullRequest struct {
	ID      int
	Team    *Team
	Commits []*Commit

	locker *sync.Mutex
}

// AssignTeam looks up for a team given a schemaName and a community
// Notice that it cleans up the current Team
func (pR *PullRequest) AssignTeam(community *Community, schName integrity.SchemaName) error {
	pR.Lock()
	pR.Team = &Team{}
	pR.Unlock()
	team, err := community.LookFor(schName)
	if err != nil {
		return err
	}
	pR.Lock()
	pR.Team = team
	pR.Unlock()
	return nil
}

func (pR *PullRequest) Lock() {
	if pR.locker == nil {
		pR.locker = new(sync.Mutex)
	}
	pR.locker.Lock()
}

func (pR *PullRequest) Unlock() {
	if pR.locker == nil {
		return
	}
	pR.locker.Unlock()
}
