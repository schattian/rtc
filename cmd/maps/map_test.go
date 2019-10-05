package maps

import (
	"sync"
	"testing"
)

func TestMap_preciseColErr(t *testing.T) {
	type fields struct {
		Name   string
		Schema []*Table
	}
	type args struct {
		colName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name:   "given column is in the map",
			fields: fields{Name: "testMap", Schema: []*Table{&Table{Name: "testTable", Columns: []string{"testCol"}}}},
			args:   args{colName: "testCol"},
			want:   errForeignColumn,
		},
		{
			name:   "given column isnt in the map",
			fields: fields{Name: "testMap", Schema: []*Table{&Table{Name: "testTable", Columns: []string{"testCol"}}}},
			args:   args{colName: "unexistantCol"},
			want:   errUnexistantColumn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Name:   tt.fields.Name,
				Schema: tt.fields.Schema,
			}
			if err := m.preciseColErr(tt.args.colName); err != tt.want {
				t.Errorf("Map.preciseColErr() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestMap_Validate(t *testing.T) {
	type fields struct {
		Name   string
		Schema []*Table
	}
	type args struct {
		tableName string
		colName   string
		wg        *sync.WaitGroup
		errCh     chan error
	}
	const maxValidationErrs = 2
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErrQt int
	}{
		{
			name:      "passes all the validations correctly",
			fields:    fields{Name: "testMap", Schema: []*Table{&Table{Name: "testTable", Columns: []string{"testCol"}}}},
			args:      args{tableName: "testTable", colName: "testCol", wg: new(sync.WaitGroup), errCh: make(chan error, maxValidationErrs)},
			wantErrQt: 0,
		},
		{
			name:      "invalid table",
			fields:    fields{Name: "testMap", Schema: []*Table{&Table{Name: "testTable", Columns: []string{"testCol"}}}},
			args:      args{tableName: "invalidTable", colName: "testCol", wg: new(sync.WaitGroup), errCh: make(chan error, maxValidationErrs)},
			wantErrQt: 1,
		},
		{
			name:      "invalid column",
			fields:    fields{Name: "testMap", Schema: []*Table{&Table{Name: "testTable", Columns: []string{"testCol"}}}},
			args:      args{tableName: "testTable", colName: "invalidCol", wg: new(sync.WaitGroup), errCh: make(chan error, maxValidationErrs)},
			wantErrQt: 1,
		},
		{
			name:      "invalid column & table",
			fields:    fields{Name: "testMap", Schema: []*Table{&Table{Name: "testTable", Columns: []string{"testCol"}}}},
			args:      args{tableName: "invalidTable", colName: "invalidCol", wg: new(sync.WaitGroup), errCh: make(chan error, maxValidationErrs)},
			wantErrQt: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Name:   tt.fields.Name,
				Schema: tt.fields.Schema,
			}
			tt.args.wg.Add(1)
			go m.Validate(tt.args.tableName, tt.args.colName, tt.args.wg, tt.args.errCh)
			tt.args.wg.Wait()
			if errQt := len(tt.args.errCh); errQt != tt.wantErrQt {
				t.Errorf("Map.Validate() obtained errQt = %v, want errQt %v", errQt, tt.wantErrQt)
			}
		})
	}
}
