package schema

import (
	"github.com/sebach1/git-crud/integrity"
)

// A Table is the representation of SQL table (or Mongo/CQL Collections) which acts as a collection of entities.
type Table struct {
	Name       integrity.TableName   `json:"name,omitempty"`
	Columns    []*Column             `json:"columns,omitempty"`
	OptionKeys []integrity.OptionKey `json:"option_keys,omitempty"`
}

func (t *Table) columnNames() (colNames []integrity.ColumnName) {
	for _, column := range t.Columns {
		colNames = append(colNames, column.Name)
	}
	return
}
