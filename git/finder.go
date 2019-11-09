package git

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sebach1/git-crud/integrity"
)

func BranchByName(ctx context.Context, DB *sql.DB, branchName integrity.BranchName) (*Branch, error) {
	branch := &Branch{}
	row := DB.QueryRowContext(ctx, fmt.Sprintf("SELECT * FROM branches WHERE Name = %s", branchName))
	err := row.Scan(branch)
	if err != nil {
		return nil, err
	}
	return branch, nil
}
