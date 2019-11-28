package schema

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/integrity"
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
			table:        gTables.Foo,
			wantColNames: []integrity.ColumnName{gColumns.Foo.Name},
		},
		{
			name:         "multi-column table",
			table:        gTables.FooBar,
			wantColNames: []integrity.ColumnName{gColumns.Foo.Name, gColumns.Bar.Name},
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
