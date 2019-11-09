package git

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sebach1/git-crud/integrity"
)

type credentials map[integrity.SchemaName]string

func (creds *credentials) Encrypt(schName integrity.SchemaName, cred string) {
}

func (creds *credentials) Decrypt(schName integrity.SchemaName, cred string) {

}

// Branch is the state-manager around indeces
type Branch struct {
	ID   int64
	Name string

	Credentials credentials

	IndexID int64
}

func NewBranch(ctx context.Context, DB *sql.DB, name integrity.BranchName) (*Branch, error) {
	res, err := DB.Exec(`INSERT INTO indeces (changes) VALUES ([])`)
	if err != nil {
		return nil, err
	}
	idxID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	branch := &Branch{Name: string(name), IndexID: idxID}
	res, err = DB.Exec(`INSERT INTO branches (name, index_id) VALUES ($1, $2)`, branch.Name, branch.IndexID)
	if err != nil {
		return nil, err
	}
	branch.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (b *Branch) Index(ctx context.Context, DB *sql.DB) (*Index, error) {
	idx := &Index{}
	row := DB.QueryRowContext(ctx, fmt.Sprintf("SELECT * FROM indeces WHERE id = %v", b.IndexID))
	err := row.Scan(idx)
	if err != nil {
		return nil, err
	}
	return idx, nil
}
