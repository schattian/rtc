package server

import (
	"github.com/sebach1/git-crud/git"
	"github.com/sebach1/git-crud/integrity"
)

type reqBody struct {
	Branch integrity.BranchName `json:"branch,omitempty"`
	Entity integrity.ID         `json:"entity,omitempty"`
	Table  integrity.TableName  `json:"table,omitempty"`
	Column integrity.ColumnName `json:"column,omitempty"`
	Value  interface{}          `json:"value,omitempty"`
	Type   integrity.CRUD       `json:"type,omitempty"`
	Opts   git.Options          `json:"opts,omitempty"`
}
