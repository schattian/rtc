package git

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/integrity"
)

// Add wraps change adding from the inferred index
func Add(
	ctx context.Context,
	db *sqlx.DB,
	entityID integrity.ID,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	branchName integrity.BranchName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) error {
	branch, err := BranchByName(ctx, db, branchName)
	if err == sql.ErrNoRows {
		branch, err = NewBranch(ctx, db, branchName)
	}
	if err != nil {
		return err
	}

	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return err
	}

	chg, err := NewChange(entityID, tableName, columnName, val, Type, opts)
	if err != nil {
		return err
	}

	err = branch.Index.Add(chg)
	if err != nil {
		return err
	}

	return nil
}

// Rm wraps change removing from the inferred index
func Rm(
	ctx context.Context,
	db *sqlx.DB,
	entityID integrity.ID,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	branchName integrity.BranchName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) error {
	branch, err := BranchByName(ctx, db, branchName)
	if err != nil {
		return err
	}

	err = branch.FetchIndex(ctx, db)
	if err != nil {
		return err
	}

	chg, err := NewChange(entityID, tableName, columnName, val, Type, opts)
	if err != nil {
		return err
	}

	err = branch.Index.Rm(chg)
	if err != nil {
		return err
	}

	return nil
}
