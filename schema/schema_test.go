package schema

import (
	"testing"

	"github.com/sebach1/git-crud/internal/integrity"
)

func TestSchema_preciseColErr(t *testing.T) {
	type args struct {
		sch     *Schema
		colName integrity.ColumnName
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "column is in the schema",
			args: args{sch: gSchemas.Basic, colName: gColumns.Basic.Name},
			want: errForeignColumn,
		},
		{
			name: "column doesn't exists in the schema",
			args: args{sch: gSchemas.Basic, colName: gColumns.Rare.Name},
			want: errNonexistentColumn,
		},
		{
			name: "schema caller is zero-valued",
			args: args{sch: gSchemas.Zero, colName: gColumns.Basic.Name},
			want: errNonexistentColumn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.sch.preciseColErr(tt.args.colName); err != tt.want {
				t.Errorf("Schema.preciseColErr() error = %v, want %v", err, tt.want)
			}
		})
	}
}
