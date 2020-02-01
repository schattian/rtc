package git

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "find branch by name")
	}
	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "index fetch")
	}
	comms, err := branch.Index.Commit(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "index commitment")
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
		return nil, errors.Wrap(err, "new change")
	}

	branch, err := BranchByName(ctx, db, branchName)
	if errors.Cause(err) == sql.ErrNoRows {
		branch, err = NewBranchWithIndex(ctx, db, branchName)
	}
	if err != nil {
		return nil, errors.Wrap(err, "find branch by name")
	}

	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "fetch index")
	}
	err = branch.Index.FetchChanges(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "fetch changes")
	}
	err = branch.Index.Add(ctx, db, chg)
	if err != nil {
		return nil, errors.Wrap(err, "index add change")
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
		return nil, errors.Wrap(err, "find branch by name")
	}

	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "index fetch")
	}
	err = branch.Index.FetchChanges(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "index changes fetch")
	}

	chg, err := NewChange(entityId, tableName, columnName, val, Type, opts)
	if err != nil {
		return nil, errors.Wrap(err, "index new change")
	}

	err = branch.Index.Rm(ctx, db, chg)
	if err != nil {
		return nil, errors.Wrap(err, "index rm change")
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
		return nil, errors.Wrap(err, "find branch by name")
	}
	commits, err := branch.UnmergedCommits(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "branch fetch unmerged commits")
	}
	pR := NewPullRequest(commits)

	own.Waiter.Add(1)
	go own.Orchestrate(ctx, community, schemaName, pR)
	err = own.WaitAndClose()
	if err != nil {
		return nil, errors.Wrap(err, "owner wait and close")
	}
	err = store.UpsertIntoDB(ctx, db, pR)
	if err != nil {
		return nil, errors.Wrap(err, "upsert pull request into db")
	}

	var comms []store.Storable
	for _, comm := range pR.Commits {
		comms = append(comms, comm)
	}
	err = store.UpsertIntoDB(ctx, db, comms...)
	if err != nil {
		return nil, errors.Wrap(err, "upsert commits into db")
	}

	return pR, nil
}
