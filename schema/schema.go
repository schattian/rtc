package schema

import (
	"errors"
	"sync"

	"github.com/sebach1/git-crud/integrity"
)

// The Schema is the representation of a Database instructive. It uses concepts of SQL.
// It provides the validation and construction structure.
type Schema struct {
	Name      integrity.SchemaName `json:"name,omitempty"`
	Blueprint []*Table             `json:"blueprint,omitempty"`
}

// Copy returns a copy of the given schema, including a deep copy if its blueprint
func (sch *Schema) Copy() *Schema {
	newSch := new(Schema)
	*newSch = *sch
	var newBlueprint []*Table
	for _, table := range newSch.Blueprint {
		newBlueprint = append(newBlueprint, table.Copy())
	}
	newSch.Blueprint = newBlueprint
	return newSch
}

// ValidateSelf performs a deep self-validation to check data integrity
// It wraps internal method validateSelf
func (sch *Schema) ValidateSelf() (err error) {
	done := make(chan bool)
	validationErrs := make(chan error)
	go sch.validateSelf(done, validationErrs)

	var errMsg string
	for {
		select {
		case <-done:
			if errMsg != "" {
				err = errors.New(errMsg)
			}
			return
		case vErr := <-validationErrs:
			errMsg += vErr.Error()
			errMsg += integrity.ErrorsSeparator
		}
	}
}

func (sch *Schema) validateSelf(done chan<- bool, vErrCh chan<- error) {
	defer func() {
		done <- true
		close(vErrCh)
	}()

	if sch == nil {
		vErrCh <- sch.validationErr(errNilSchema)
		return
	}

	tablesQt := len(sch.Blueprint)
	if tablesQt == 0 {
		vErrCh <- sch.validationErr(errNilBlueprint)
	}

	var schVWg sync.WaitGroup
	schVWg.Add(tablesQt)
	for _, table := range sch.Blueprint {
		go table.validateSelf(&schVWg, vErrCh)
	}

	if sch.Name == "" {
		vErrCh <- sch.validationErr(errNilSchemaName)
	}

	schVWg.Wait()
}

func (sch *Schema) validationErr(err error) *integrity.ValidationError {
	var name string
	if sch == nil {
		name = ""
	} else {
		name = string(sch.Name)
	}
	return &integrity.ValidationError{Err: err, Origin: "schema", OriginName: name}
}

// ValidateCtx checks if the context of the given tableName and colName is valid
// Notice that, as well as the wrapper validations should provoke a chained
// of undesired (and maybe more confusing than clear) errs, the errCh should be buffered w/sz=1
func (sch *Schema) ValidateCtx(
	tableName integrity.TableName,
	colName integrity.ColumnName,
	optionKeys []integrity.OptionKey,
	val interface{},
	helperScope *Planisphere,
	wg *sync.WaitGroup,
	errCh chan<- error,
) {
	defer wg.Done()

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

	if colName == "" {
		return
	}

	for _, col := range table.Columns {
		if colName == col.Name {
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

// Wraps Column.applyBuiltinValidator() over all cols
func (sch *Schema) applyBuiltinValidators() (err error) {
	for _, table := range sch.Blueprint {
		for _, col := range table.Columns {
			err = col.applyBuiltinValidator()
			if err != nil {
				return
			}
		}
	}
	return nil
}
