package git

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/integrity"
)

func BranchByName(ctx context.Context, db *sqlx.DB, branchName integrity.BranchName) (*Branch, error) {
	branch := Branch{}
	err := db.GetContext(ctx, &branch, `SELECT * FROM branches WHERE name=?`, branchName)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}
