package schema

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTable_columnNames(t *testing.T) {
	tests := []struct {
		name         string
		table        *Table
		wantColNames []ColumnName
	}{
		{
			name:         "single-column table",
			table:        gTables.Basic,
			wantColNames: []ColumnName{gColumns.Basic.Name},
		},
		{
			name:         "multi-column table",
			table:        gTables.BasicRare,
			wantColNames: []ColumnName{gColumns.Basic.Name, gColumns.Rare.Name},
		},
		{
			name:         "table is zero-valued",
			table:        gTables.Zero,
			wantColNames: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotColNames := tt.table.columnNames(); !reflect.DeepEqual(gotColNames, tt.wantColNames) {
				t.Errorf(cmp.Diff(gotColNames, tt.wantColNames))
				t.Errorf("Table.columnNames() = %v, want %v", gotColNames, tt.wantColNames)
			}
		})
	}
}
