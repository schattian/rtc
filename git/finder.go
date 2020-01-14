package git

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/rtc/integrity"
)

// BranchByName finds a branch in the DB given its name
func BranchByName(ctx context.Context, db *sqlx.DB, branchName integrity.BranchName) (*Branch, error) {
	branch := Branch{}
	err := db.GetContext(ctx, &branch, `SELECT * FROM branches WHERE name=?`, branchName)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

// CommitById finds a commit in the DB given its id
func CommitById(ctx context.Context, db *sqlx.DB, id int64) (*Commit, error) {
	comm := &Commit{}
	err := db.GetContext(ctx, comm, `SELECT * FROM commits WHERE id=?`, id)
	if err != nil {
		return nil, err
	}
	return comm, nil
}
