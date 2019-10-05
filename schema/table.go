package schema

func (t *Table) columnNames() (colNames []ColumnName) {
	if t == nil {
		return
	}
	for _, column := range t.Columns {
		colNames = append(colNames, column.Name)
	}
	return
}
