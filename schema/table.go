package schema

// A Table is the representation of SQL table (or Mongo/CQL Collections) which acts as a collection of entities.
type Table struct {
	Name    TableName `json:"name,omitempty"`
	Columns []*Column `json:"columns,omitempty"`
	// Maintainer git.Collaborator `json:"maintainer,omitempty"`
}

// TableName is the name of a table
type TableName string

func (t *Table) columnNames() (colNames []ColumnName) {
	for _, column := range t.Columns {
		colNames = append(colNames, column.Name)
	}
	return
}
