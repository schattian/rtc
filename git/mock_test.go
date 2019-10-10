package git

import (
	"context"

	"github.com/sebach1/git-crud/internal/integrity"
)

type collabMock struct {
	Err error
}

func (mock *collabMock) Push(ctx context.Context, comm *Commit) (*Commit, error) {
	if err := mock.Err; err != nil {
		return nil, err
	}
	return comm, nil
}

func (mock *collabMock) Pull(ctx context.Context, comm *Commit) (*Commit, error) {
	if err := mock.Err; err != nil {
		return nil, err
	}
	return comm, nil
}

func (mock *collabMock) Delete(ctx context.Context, comm *Commit) (*Commit, error) {
	if err := mock.Err; err != nil {
		return nil, err
	}
	return comm, nil
}

func (mock *collabMock) Init(ctx context.Context) error {
	if err := mock.Err; err != nil {
		return err
	}
	return nil
}

func (chg *Change) changeType(newType integrity.CRUD) *Change {
	chg.Type = newType
	return chg
}

// Pushes a MOCKED COLLABORATOR with the ASSIGNED TABLE which RETURNS THE GIVEN ERROR
func (pR *PullRequest) mock(tableName integrity.TableName, err error) *PullRequest {
	pR.Team.mock(tableName, err)
	return pR
}

// Copies and pushes a MOCKED COLLABORATOR with the ASSIGNED TABLE which RETURNS THE GIVEN ERROR
func (t *Team) mock(tableName integrity.TableName, err error) *Team {
	mock := &collabMock{Err: err}
	t.Members = append(t.Members, &Member{AssignedTable: tableName, Collab: mock})
	return t
}

func (pR *PullRequest) addCommit(comm *Commit) *PullRequest {
	pullRequest := pR.copy()
	pullRequest.Commits = append(pullRequest.Commits, comm)
	return pullRequest
}

func (m *Member) copy() *Member {
	member := new(Member)
	*member = *m
	return member
}

func (t *Team) copy() *Team {
	team := new(Team)
	*team = *t
	var newMembers []*Member
	for _, member := range t.Members {
		newMembers = append(newMembers, member)
	}
	team.Members = newMembers
	return team
}

func (chg *Change) copy() *Change {
	newChg := new(Change)
	*newChg = *chg
	return newChg
}

func (comm *Commit) copy() *Commit {
	newComm := new(Commit)
	*newComm = *comm
	var newChgs []*Change
	for _, chg := range comm.Changes {
		newChgs = append(newChgs, chg.copy())
	}
	newComm.Changes = newChgs
	return newComm
}

func (pR *PullRequest) copy() *PullRequest {
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
