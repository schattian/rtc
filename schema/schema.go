package schema

import (
	"errors"
	"sync"

	"github.com/sebach1/git-crud/integrity"
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
	optionKeys []integrity.OptionKey,
	val interface{},
	helperScope *Planisphere,
	wg *sync.WaitGroup,
	errCh chan<- error,
) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			errCh <- errors.New("schema.Validate() unhandled PANIC")
		}
	}()

	table, err := sch.tableByName(tableName, helperScope)
	if err != nil {
		errCh <- err
		return
	}

	for _, key := range optionKeys {
		if !table.optionKeyIsValid(key) {
			errCh <- errInvalidOptionKey
			return
		}
	}

	if colName == "" { // Skip column validation (before the change MUST BE TYPE VALIDATED)
		return
	}

	for _, col := range table.Columns {
		if colName == col.Name {
			if val == nil {
				return
			}
			err = col.Validate(val)
			if err != nil {
				errCh <- err
				return
			}
			return
		}
	}
	errCh <- sch.preciseColErr(colName)
}

func (t *Table) optionKeyIsValid(key integrity.OptionKey) bool {
	for _, validKey := range t.OptionKeys {
		if validKey == key {
			return true
		}
	}
	return false
}

func (sch *Schema) tableByName(tableName integrity.TableName, helperScope *Planisphere) (*Table, error) {
	if len(sch.Blueprint) == 0 {
		return nil, errNilBlueprint
	}
	for _, table := range sch.Blueprint {
		if tableName == table.Name {
			return table, nil
		}
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
