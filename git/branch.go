package git

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/rtc/integrity"
	"github.com/sebach1/rtc/internal/store"
)

// Branch is the state-manager around indices
type Branch struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	Credentials credentials

	Index *Index `json:"index,omitempty"`

	IndexId int64 `json:"index_id,omitempty"`
}

func (b *Branch) UnmergedCommits(ctx context.Context, db *sqlx.DB) ([]*Commit, error) {
	rows, err := db.NamedQueryContext(ctx, `SELECT * FROM commits WHERE merged=false AND branch_id=:id`, b)
	if err != nil {
		return nil, err
	}
	var comms []*Commit
	defer rows.Close()
	for rows.Next() {
		comm := &Commit{}
		err = rows.StructScan(comm)
		if err != nil {
			return nil, err
		}
		comms = append(comms, comm)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comms, nil
}

// NewBranchWithIndex safety creates a new Branch entity and assigns a new index_id to it
// Notice it persists on the db and assigns the inserted id
func NewBranchWithIndex(ctx context.Context, db *sqlx.DB, name integrity.BranchName) (*Branch, error) {
	if ctx == nil {
		ctx = context.Background() // Avoid panicking on .ExecContext()
	}
	res, err := db.ExecContext(ctx, `INSERT INTO indices DEFAULT VALUES`)
	if err != nil {
		return nil, err
	}
	idxId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	branch := &Branch{Name: string(name), IndexId: idxId}
	err = store.InsertIntoDB(ctx, db, branch)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

// FetchIndex retrieves the Index by .IndexId and assigns it to .Index field
func (b *Branch) FetchIndex(ctx context.Context, db *sqlx.DB) error {
	if b.IndexId == 0 {
		return errNilIndexId
	}
	row := db.QueryRowxContext(ctx, `SELECT * FROM indices WHERE id=?`, b.IndexId)
	if b.Index == nil { // Avoid panicking on .StructScan due nil receiver
		b.Index = &Index{}
	}
	err := row.StructScan(b.Index)
	if err != nil {
		return err
	}
	return nil
}
