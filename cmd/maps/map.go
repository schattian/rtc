package maps

import (
	"sync"
)

// The Map is the representation of a Database instructive. It uses concepts of SQL.
// The Schema provided by the map gives the validation structure.
type Map struct {
	Name   string
	Schema []*Table
}

// colNames plucks the columns from its tables
func (m *Map) colNames() (colNames []string) {
	for _, table := range m.Schema {
		for _, column := range table.Columns {
			colNames = append(colNames, column)
		}
	}
	return
}

// tableNames plucks the name from its tables
func (m *Map) tableNames() (tableNames []string) {
	for _, table := range m.Schema {
		tableNames = append(tableNames, table.Name)
	}
	return
}

// colsByTableName returns the column names given the parent' table name
func (m *Map) colsByTableName(tableName string) ([]string, error) {
	for _, table := range m.Schema {
		if tableName == table.Name {
			return table.Columns, nil
		}
	}
	return nil, preciseTableErr(Planisphere, tableName)
}

func (m *Map) validateTable(tableName string, wg *sync.WaitGroup, errCh chan<- error) {
	cols, err := m.colsByTableName(tableName)
	if err != nil {
		errCh <- err
	}
}

// Validate checks if the context of the given tableName and colName is valid
func (m *Map) Validate(tableName, colName string, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	cols, err := m.colsByTableName(tableName)
	if err != nil {
		errCh <- err
		return
	}

	for _, col := range cols {
		if colName == col {
			return
		}
	}
	errCh <- m.preciseColErr(colName)
}

// preciseColErr gives a more accurated error to a validation of a column
// It assumes the column is errored, and checks if it exists or if instead its a context err
func (m *Map) preciseColErr(colName string) (err error) {
	for _, column := range m.colNames() {
		if column == colName {
			return errForeignColumn
		}
	}
	return errUnexistantColumn
}
