package schema

import (
	"sync"

	"github.com/sebach1/git-crud/integrity"
)

// A Table is the representation of SQL table (or Mongo/CQL Collections) which acts as a collection of entities.
type Table struct {
	Name       integrity.TableName   `json:"name,omitempty"`
	Columns    []*Column             `json:"columns,omitempty"`
	OptionKeys []integrity.OptionKey `json:"option_keys,omitempty"`
}

func (t *Table) validateSelf(wg *sync.WaitGroup, vErrCh chan<- error) {
	defer func() {
		wg.Done()
	}()

	if t == nil {
		vErrCh <- t.validationErr(errNilTable)
		return
	}

	colsQt := len(t.Columns)
	if colsQt == 0 {
		vErrCh <- t.validationErr(errNilColumns)
	}

	var tVWg sync.WaitGroup
	tVWg.Add(colsQt)
	for _, col := range t.Columns {
		go col.validateSelf(&tVWg, vErrCh)
	}

	if t.Name == "" {
		vErrCh <- t.validationErr(errNilTableName)
	}

	tVWg.Wait()
}

func (t *Table) validationErr(err error) *integrity.ValidationError {
	var name string
	if t == nil {
		name = ""
	} else {
		name = string(t.Name)
	}
	return &integrity.ValidationError{Err: err, OriginType: "table", OriginName: name}
}

// Copy returns a copy of the given table, including embedded cols
func (t *Table) Copy() *Table {
	newTab := new(Table)
	*newTab = *t
	var newCols []*Column
	for _, col := range newTab.Columns {
		newCols = append(newCols, col.Copy())
	}
	newTab.Columns = newCols
	return newTab
}

func (t *Table) columnNames() (colNames []integrity.ColumnName) {
	for _, column := range t.Columns {
		colNames = append(colNames, column.Name)
	}
	return
}
