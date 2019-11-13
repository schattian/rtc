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

	Index *Index

	IndexID int64
}

// NewBranch safety creates a new Branch entity
// Notice it doesn't saves it on the db
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

// FetchIndex retrieves the Index by .IndexID and assigns it to .Index field
func (b *Branch) FetchIndex(ctx context.Context, db *sqlx.DB) error {
	row := db.QueryRowxContext(ctx, `SELECT * FROM indeces WHERE id=?`, b.IndexID)
	err := row.StructScan(b.Index)
	if err != nil {
		return err
	}
	return nil
}
