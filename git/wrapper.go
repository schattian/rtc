package git

import (
	"context"
	"database/sql"

	"github.com/sebach1/git-crud/integrity"
)

func Add(
	ctx context.Context,
	DB *sql.DB,
	entityID integrity.ID,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	branchName integrity.BranchName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) error {
	branch, err := BranchByName(ctx, DB, branchName)
	if err == sql.ErrNoRows {
		branch, err = NewBranch(ctx, DB, branchName)
	}
	if err != nil {
		return err
	}

	idx, err := branch.Index(ctx, DB)
	if err != nil {
		return err
	}

	chg, err := NewChange(entityID, tableName, columnName, val, Type, opts)
	if err != nil {
		return err
	}

	err = idx.Add(chg)
	if err != nil {
		return err
	}

	return nil
}

func Rm(
	ctx context.Context,
	DB *sql.DB,
	entityID integrity.ID,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	branchName integrity.BranchName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) error {
	branch, err := BranchByName(ctx, DB, branchName)
	if err != nil {
		return err
	}

	idx, err := branch.Index(ctx, DB)
	if err != nil {
		return err
	}

	chg, err := NewChange(entityID, tableName, columnName, val, Type, opts)
	if err != nil {
		return err
	}

	err = idx.Rm(chg)
	if err != nil {
		return err
	}

	return nil
}
