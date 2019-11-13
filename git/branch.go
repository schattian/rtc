package git

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/integrity"
)

// Branch is the state-manager around indeces
type Branch struct {
	ID   int64
	Name string

	Credentials credentials

	IndexID int64
}

func NewBranch(ctx context.Context, db *sqlx.DB, name integrity.BranchName) (*Branch, error) {
	res, err := db.Exec(`INSERT INTO indeces (changes) VALUES ([])`)
	if err != nil {
		return nil, err
	}
	idxID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	branch := &Branch{Name: string(name), IndexID: idxID}
	res, err = db.Exec(`INSERT INTO branches (name, index_id) VALUES ($1, $2)`, branch.Name, branch.IndexID)
	if err != nil {
		return nil, err
	}
	branch.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (b *Branch) Index(ctx context.Context, db *sqlx.DB) (*Index, error) {
	idx := &Index{}
	row := db.QueryRowxContext(ctx, `SELECT * FROM indeces WHERE id=?`, b.IndexID)
	err := row.StructScan(idx)
	if err != nil {
		return nil, err
	}
	return idx, nil
}
