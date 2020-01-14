package git

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/rtc/integrity"
	"github.com/sebach1/rtc/internal/store"
	"github.com/sebach1/rtc/schema"
)

func Comm(
	ctx context.Context,
	db *sqlx.DB,
	branchName integrity.BranchName,
) ([]*Commit, error) {
	branch, err := BranchByName(ctx, db, branchName)
	if err != nil {
		return nil, err
	}
	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return nil, err
	}
	comms, err := branch.Index.Commit(ctx, db)
	if err != nil {
		return nil, err
	}
	return comms, nil
}

// Add wraps change adding from the inferred index
func Add(
	ctx context.Context,
	db *sqlx.DB,
	entityId integrity.Id,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	branchName integrity.BranchName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) (*Change, error) {
	chg, err := NewChange(entityId, tableName, columnName, val, Type, opts) // The specific order is to avoid creating new branch with unvalid change
	if err != nil {
		return nil, err
	}

	branch, err := BranchByName(ctx, db, branchName)
	if err == sql.ErrNoRows {
		branch, err = NewBranchWithIndex(ctx, db, branchName)
	}
	if err != nil {
		return nil, err
	}

	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return nil, err
	}
	err = branch.Index.FetchChanges(ctx, db)
	if err != nil {
		return nil, err
	}
	err = branch.Index.Add(ctx, db, chg)
	if err != nil {
		return nil, err
	}
	return chg, nil
}

// Rm wraps change removal from the inferred index
func Rm(
	ctx context.Context,
	db *sqlx.DB,
	entityId integrity.Id,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	branchName integrity.BranchName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) (*Change, error) {
	branch, err := BranchByName(ctx, db, branchName)
	if err != nil {
		return nil, err
	}

	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return nil, err
	}
	err = branch.Index.FetchChanges(ctx, db)
	if err != nil {
		return nil, err
	}

	chg, err := NewChange(entityId, tableName, columnName, val, Type, opts)
	if err != nil {
		return nil, err
	}

	err = branch.Index.Rm(ctx, db, chg)
	if err != nil {
		return nil, err
	}
	return chg, nil
}

func Orchestrate(
	ctx context.Context,
	db *sqlx.DB,
	project *schema.Planisphere,
	branchName integrity.BranchName,
	schemaName integrity.SchemaName,
	community *Community,
) (*PullRequest, error) {
	own, err := NewOwner(project)
	if err != nil {
		return nil, err
	}
	branch, err := BranchByName(ctx, db, branchName)
	if err != nil {
		return nil, err
	}
	commits, err := branch.UnmergedCommits(ctx, db)
	if err != nil {
		return nil, err
	}
	pR := NewPullRequest(commits)

	own.Waiter.Add(1)
	go own.Orchestrate(ctx, community, schemaName, pR)
	err = own.WaitAndClose()
	if err != nil {
		return nil, err
	}
	err = store.UpsertIntoDB(ctx, db, pR)
	if err != nil {
		return nil, err
	}

	var comms []store.Storable
	for _, comm := range pR.Commits {
		comms = append(comms, comm)
	}
	err = store.UpsertIntoDB(ctx, db, comms...)
	if err != nil {
		return nil, err
	}

	return pR, nil
}
