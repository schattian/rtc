package schema

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/integrity"
)

func TestTable_columnNames(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		table        *Table
		wantColNames []integrity.ColumnName
	}{
		{
			name:         "single-column table",
			table:        gTables.Basic,
			wantColNames: []integrity.ColumnName{gColumns.Basic.Name},
		},
		{
			name:         "multi-column table",
			table:        gTables.BasicRare,
			wantColNames: []integrity.ColumnName{gColumns.Basic.Name, gColumns.Rare.Name},
		},
		{
			name:         "table is zero-valued",
			table:        gTables.Zero,
			wantColNames: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if gotColNames := tt.table.columnNames(); !reflect.DeepEqual(gotColNames, tt.wantColNames) {
				t.Errorf(cmp.Diff(gotColNames, tt.wantColNames))
				t.Errorf("Table.columnNames() = %v, want %v", gotColNames, tt.wantColNames)
			}
		})
	}
}
