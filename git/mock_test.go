package git

import (
	"context"
	"log"

	"github.com/sebach1/rtc/integrity"
)

type collabMock struct {
	Err error
}

func (mock *collabMock) Create(ctx context.Context, comm *Commit) (*Commit, error) {
	if err := mock.Err; err != nil {
		return nil, err
	}
	return comm, nil
}

func (mock *collabMock) Retrieve(ctx context.Context, comm *Commit) (*Commit, error) {
	if err := mock.Err; err != nil {
		return nil, err
	}
	return comm, nil
}

func (mock *collabMock) Update(ctx context.Context, comm *Commit) (*Commit, error) {
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
