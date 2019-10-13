package schema

import (
	"sync"
	"testing"

	"github.com/sebach1/git-crud/internal/integrity"
)

func TestSchema_preciseColErr(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.args.sch.preciseColErr(tt.args.colName); err != tt.want {
				t.Errorf("Schema.preciseColErr() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSchema_Validate(t *testing.T) {
	t.Parallel()
	type args struct {
		tableName   integrity.TableName
		colName     integrity.ColumnName
		optionKeys  []integrity.OptionKey
		helperScope *Planisphere
	}
	tests := []struct {
		name    string
		sch     *Schema
		args    args
		wantErr bool
	}{
		{
			name: "passes the validations",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Basic.Name,
				optionKeys:  []integrity.OptionKey{gTables.Basic.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: false,
		},
		{
			name: "optionKey nonexistant",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Basic.Name,
				optionKeys:  []integrity.OptionKey{gTables.Rare.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
		{
			name: "column not inside any table",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Rare.Name,
				optionKeys:  []integrity.OptionKey{gTables.Basic.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
		{
			name: "table nonexistant",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Rare.Name,
				colName:     gColumns.Basic.Name,
				optionKeys:  []integrity.OptionKey{},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wg := new(sync.WaitGroup)
			errCh := make(chan error, 1)
			wg.Add(1)
			go tt.sch.Validate(tt.args.tableName, tt.args.colName, tt.args.optionKeys, tt.args.helperScope, wg, errCh)
			wg.Wait()
			isErrored := len(errCh) == 1
			if isErrored && !tt.wantErr {
				err := <-errCh
				t.Errorf("Schema.Validate() error: %v; wantErr %v", err, tt.wantErr)
			}
			if !isErrored && tt.wantErr {
				t.Errorf("Schema.Validate() error: %v; wantErr %v", nil, tt.wantErr)
			}
		})
	}
}
