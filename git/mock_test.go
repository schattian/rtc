package git

import (
	"context"

	"github.com/sebach1/git-crud/internal/integrity"
)

type collabMock struct {
	Err error
}

func (t *Team) copy() *Team {
	team := new(Team)
	*team = *t
	return team
}

func (chg *Change) copy() *Change {
	newChg := new(Change)
	*newChg = *chg
	return newChg
}

func (chg *Change) changeType(newType integrity.CRUD) *Change {
	chg.Type = newType
	return chg
}

func (pR *PullRequest) copy() *PullRequest {
	newPr := new(PullRequest)
	*newPr = *pR
	return newPr
}

func (t *Team) mockedCopy(tableName integrity.TableName, err error) *Team {
	team := t.copy()
	mock := &collabMock{Err: err}
	team.Members = append(team.Members, &Member{AssignedTable: tableName, Collab: mock})
	return team
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
