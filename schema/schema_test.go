package schema

import (
	"sync"
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

func TestSchema_Validate(t *testing.T) {
	type args struct {
		tableName   integrity.TableName
		colName     integrity.ColumnName
		helperScope *Planisphere
		wg          *sync.WaitGroup
		errCh       chan error
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
				helperScope: &Planisphere{gSchemas.Basic},
				wg:          new(sync.WaitGroup),
				errCh:       make(chan error, 1),
			},
			wantErr: false,
		},
		{
			name: "column not inside any table",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Rare.Name,
				helperScope: &Planisphere{gSchemas.Basic},
				wg:          new(sync.WaitGroup),
				errCh:       make(chan error, 1),
			},
			wantErr: true,
		},
		{
			name: "table nonexistant",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Rare.Name,
				colName:     gColumns.Basic.Name,
				helperScope: &Planisphere{gSchemas.Basic},
				wg:          new(sync.WaitGroup),
				errCh:       make(chan error, 1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.wg.Add(1)
			go tt.sch.Validate(tt.args.tableName, tt.args.colName, tt.args.helperScope, tt.args.wg, tt.args.errCh)
			tt.args.wg.Wait()
			isErrored := len(tt.args.errCh) == 1
			if isErrored && !tt.wantErr {
				err := <-tt.args.errCh
				t.Errorf("Schema.Validate() error: %v; wantErr %v", err, tt.wantErr)
			}
			if !isErrored && tt.wantErr {
				t.Errorf("Schema.Validate() error: %v; wantErr %v", nil, tt.wantErr)
			}
		})
	}
}
