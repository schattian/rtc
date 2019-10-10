package schema

import (
	"sync"

	"github.com/sebach1/git-crud/internal/integrity"
)

// The Schema is the representation of a Database instructive. It uses concepts of SQL.
// The Schema provided by the schema gives the validation structure.
type Schema struct {
	Name      integrity.SchemaName `json:"name,omitempty"`
	Blueprint []*Table             `json:"blueprint,omitempty"`
}

// Validate checks if the context of the given tableName and colName is valid
// Notice that, as well as the wrapper validations should provoke a chained
// of undesired (and maybe more confusing than clear) errs, the errCh should be buffered w/sz=1
func (sch *Schema) Validate(
	tableName integrity.TableName,
	colName integrity.ColumnName,
	helperScope *Planisphere,
	wg *sync.WaitGroup,
	errCh chan<- error,
) {
	defer wg.Done()

	cols, err := sch.colsByTableName(tableName, helperScope)
	if err != nil {
		errCh <- err
		return
	}

	for _, col := range cols {
		if colName == col {
			return
		}
	}
	errCh <- sch.preciseColErr(colName)
}

// colsByTableName returns the column names given the parent' table name
func (sch *Schema) colsByTableName(tableName integrity.TableName, helperScope *Planisphere) ([]integrity.ColumnName, error) {
	for _, table := range sch.Blueprint {
		if tableName != table.Name {
			continue
		}
		return table.columnNames(), nil
	}
	return nil, helperScope.preciseTableErr(tableName)
}

// colNames plucks all the columnNames from its tables
func (sch *Schema) colNames() (colNames []integrity.ColumnName) {
	for _, table := range sch.Blueprint {
		for _, column := range table.Columns {
			colNames = append(colNames, column.Name)
		}
	}
	return
}

// tableNames plucks the name from its tables
func (sch *Schema) tableNames() (tableNames []integrity.TableName) {
	for _, table := range sch.Blueprint {
		tableNames = append(tableNames, table.Name)
	}
	return
}

// preciseColErr gives a more accurate error to a validation of a column
// It assumes the column is errored, and checks if it exists or if instead its a context err
func (sch *Schema) preciseColErr(colName integrity.ColumnName) (err error) {
	for _, column := range sch.colNames() {
		if column == colName {
			return errForeignColumn
		}
	}
	return errNonexistentColumn
}
