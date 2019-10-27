package fabric

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/schema"
	"github.com/spf13/afero"
)

func TestFabric_Produce(t *testing.T) {
	t.Parallel()
	type args struct {
		marshal string
	}
	customDir := "testF/testSf"
	tests := []struct {
		name    string
		fabric  *Fabric
		args    args
		wantDir string
		wantErr bool
	}{
		{
			name:    "correct usage",
			fabric:  &Fabric{Schema: gSchemas.Basic},
			args:    args{marshal: "json"},
			wantDir: fmt.Sprintf("fabric/%v", strings.ToLower(string(gSchemas.Basic.Name))),
			wantErr: false,
		},
		{
			name:    "but w/PRESET DIR",
			fabric:  &Fabric{Schema: gSchemas.Basic, Dir: customDir},
			args:    args{marshal: "json"},
			wantDir: customDir,
			wantErr: false,
		},
		{
			name:    "SCHEMA is NIL",
			fabric:  &Fabric{Dir: customDir}, // customDir: see that checking os existance of "" dir always returns true
			args:    args{marshal: "json"},
			wantDir: customDir,
			wantErr: true,
		},
		{
			name:    "SCHEMA NAME is NIL",
			fabric:  &Fabric{Schema: gSchemas.Zero, Dir: customDir}, // customDir: see that checking os existance of "" dir always returns true
			args:    args{marshal: "json"},
			wantDir: customDir,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mFs := afero.NewMemMapFs()
			err := tt.fabric.Produce(tt.args.marshal, mFs)
			isCreated, osErr := afero.DirExists(mFs, tt.fabric.Dir)
			if osErr != nil {
				t.Errorf("Fabric.Produce() got UNEXPECTED ERR when trying to fetch dir on aferos' MemMapFs: %v", err)
			}
			if tt.fabric.Dir != tt.wantDir {
				t.Errorf("Fabric.Produce() DIFF DIR than expected; want: %v, got: %v", tt.wantDir, tt.fabric.Dir)
			}

			if err != nil {
				if isCreated {
					t.Error("Fabric.Produce() GENERATED the DIRectory. EXPECTed to NOT generate it")
				}
				return
			}

			if !isCreated {
				t.Error("Fabric.Produce() did NOT GENERATE the DIRectory. EXPECTed TO generate it")
			}
		})
	}
}

// func TestFabric_structFromTable(t *testing.T) {
// 	type fields struct {
// 		Schema *schema.Schema
// 		Dir    string
// 		wg     *sync.WaitGroup
// 		fsWg   *sync.WaitGroup
// 		fsSmp  chan int
// 	}
// 	type args struct {
// 		table   *schema.Table
// 		marshal string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *tableData
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Fabric{
// 				Schema: tt.fields.Schema,
// 				Dir:    tt.fields.Dir,
// 				wg:     tt.fields.wg,
// 				fsWg:   tt.fields.fsWg,
// 				fsSmp:  tt.fields.fsSmp,
// 			}
// 			if got := f.structFromTable(tt.args.table, tt.args.marshal); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Fabric.structFromTable() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// Locktest
func Test_fieldFromColumn(t *testing.T) {
	type args struct {
		col *schema.Column
	}
	// rmType := func(col *schema.Column) *schema.Column { col.Type = ""; return col }
	// rmName := func(col *schema.Column) *schema.Column { col.Name = ""; return col }
	tests := []struct {
		name    string
		args    args
		want    *columnData
		wantErr bool
	}{
		{
			name: "BASIC COLUMN conversion",
			args: args{col: gColumns.Basic},
			want: &columnData{
				Name: toCamelCase(string(gColumns.Basic.Name)),
				Type: gColumns.Basic.Type,
				Tag:  toSnakeCase(string(gColumns.Basic.Name), '_'),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fieldFromColumn(tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("fieldFromColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("fieldFromColumn() mismatch (-want +got): %s", diff)
			}
		})
	}
}
