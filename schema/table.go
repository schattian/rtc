package schema

import (
	"sync"

	"github.com/sebach1/rtc/integrity"
	"github.com/sebach1/rtc/internal/xerrors"
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

func (t *Table) validationErr(err error) *xerrors.ValidationError {
	var name string
	if t == nil {
		name = ""
	} else {
		name = string(t.Name)
	}
	return &xerrors.ValidationError{Err: err, OriginType: "table", OriginName: name}
}

func (t *Table) columnNames() (colNames []integrity.ColumnName) {
	for _, column := range t.Columns {
		colNames = append(colNames, column.Name)
	}
	return
}
