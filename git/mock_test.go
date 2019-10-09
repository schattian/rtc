package git

import (
	"context"

	"github.com/sebach1/git-crud/internal/integrity"
)

type collabMock struct {
	Err error
}

func (t *Team) mockedCopy(tableName integrity.TableName, err error) *Team {
	team := new(Team)
	*team = *t
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
