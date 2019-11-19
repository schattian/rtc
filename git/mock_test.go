package git

import (
	"context"
	"fmt"
	"log"

	"github.com/sebach1/git-crud/integrity"
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
	mockingErr := t.AddMember(tableName, mock, true /*forces mock*/)
	if mockingErr != nil {
		log.Fatal(mockingErr)
	}
	return t
}

func (pR *PullRequest) addCommit(comm *Commit) *PullRequest {
	pullRequest := pR.copy()
	pullRequest.Commits = append(pullRequest.Commits, comm)
	return pullRequest
}

func (opts *Options) assignAndReturn(k integrity.OptionKey, v interface{}) Options {
	newOpts := *opts
	newOpts[k] = v
	*opts = newOpts
	return *opts
}

func (m *Member) copy() *Member {
	panicIfNilAtCopy(m, "MEMBER")
	member := new(Member)
	*member = *m
	return member
}

func (t *Team) copy() *Team {
	panicIfNilAtCopy(t, "TEAM")
	team := new(Team)
	*team = *t
	var newMembers []*Member
	for _, member := range t.Members {
		newMembers = append(newMembers, member.copy())
	}
	team.Members = newMembers
	return team
}

func (opts *Options) copy() Options {
	panicIfNilAtCopy(opts, "OPTIONS")
	newOpts := make(Options)
	for k, v := range *opts {
		newOpts[k] = v
	}
	return newOpts
}

func (chg *Change) copy() *Change {
	panicIfNilAtCopy(chg, "CHANGE")
	newChg := &Change{}
	*newChg = *chg
	if chg.Options != nil {
		newChg.Options = chg.Options.copy()
	}
	return newChg
}

func panicIfNilAtCopy(elem interface{}, msg string) {
	if elem == nil {
		panic(fmt.Sprintf("the %v must be != nil to execute the copy", msg))
	}
}

func (comm *Commit) copy() *Commit {
	panicIfNilAtCopy(comm, "COMMIT")
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
	panicIfNilAtCopy(pR, "PULL REQUEST")
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
