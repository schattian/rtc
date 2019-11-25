package git

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/internal/store"
)

// Branch is the state-manager around indeces
type Branch struct {
	Id   int64
	Name string

	Credentials credentials

	Index *Index

	IndexId int64
}

// NewBranchWithIndex safety creates a new Branch entity and assigns a new index_id to it
// Notice it persists on the db and assigns the inserted id
func NewBranchWithIndex(ctx context.Context, db *sqlx.DB, name integrity.BranchName) (*Branch, error) {
	res, err := db.Exec(`INSERT INTO indeces DEFAULT VALUES`)
	if err != nil {
		return nil, err
	}
	idxId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	branch := &Branch{Name: string(name), IndexId: idxId}
	id, err := store.InsertToDB(ctx, branch, db)
	if err != nil {
		return nil, err
	}
	branch.SetId(id)
	return branch, nil
}

// FetchIndex retrieves the Index by .IndexId and assigns it to .Index field
func (b *Branch) FetchIndex(ctx context.Context, db *sqlx.DB) error {
	row := db.QueryRowxContext(ctx, `SELECT * FROM indeces WHERE id=?`, b.IndexId)
	err := row.StructScan(b.Index)
	if err != nil {
		return err
	}
	return nil
}
