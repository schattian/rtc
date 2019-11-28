package schema

import (
	"reflect"
	"testing"

	"github.com/sebach1/rtc/integrity"
)

func Test_preciseTableErr(t *testing.T) {
	type args struct {
		tableName integrity.TableName
	}

	tests := []struct {
		name    string
		args    args
		psph    Planisphere
		wantErr error
	}{
		{
			name:    "given tableName is in a schema",
			args:    args{tableName: gTables.Foo.Name},
			psph:    Planisphere{gSchemas.Foo},
			wantErr: errForeignTable,
		},
		{
			name:    "given tableName doesn't exists on any scoped schema",
			args:    args{gTables.Zero.Name},
			psph:    Planisphere{gSchemas.Foo, gSchemas.Bar},
			wantErr: errNonexistentTable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.psph.preciseTableErr(tt.args.tableName); err != tt.wantErr {
				t.Errorf("preciseTableErr() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanisphere_GetSchemaFromName(t *testing.T) {
	t.Parallel()
	type args struct {
		schemaName integrity.SchemaName
	}

	tests := []struct {
		name    string
		psph    Planisphere
		args    args
		want    *Schema
		wantErr error
	}{
		{
			name:    "giving a single schema matching the given name",
			args:    args{gSchemas.Foo.Name},
			psph:    Planisphere{gSchemas.Foo},
			want:    gSchemas.Foo,
			wantErr: nil,
		},
		{
			name:    "giving multiple gSchemas matching the given name it matches the first",
			args:    args{gSchemas.Foo.Name},
			psph:    Planisphere{gSchemas.Foo, gSchemas.FooBar},
			want:    gSchemas.Foo,
			wantErr: nil,
		},
		{
			name:    "giving mixed correct schemas w/nil",
			args:    args{gSchemas.Foo.Name},
			psph:    Planisphere{nil, gSchemas.Foo, nil},
			want:    gSchemas.Foo,
			wantErr: nil,
		},
		//
		{
			name:    "blank schemaName and having blank schemas",
			args:    args{gSchemas.Zero.Name},
			psph:    Planisphere{gSchemas.Zero, gSchemas.Zero},
			want:    nil,
			wantErr: errSchemaNotFoundInScope,
		},
		{
			name:    "giving no schema",
			args:    args{gSchemas.Foo.Name},
			psph:    Planisphere{},
			want:    nil,
			wantErr: errSchemaNotFoundInScope,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.psph.GetSchemaFromName(tt.args.schemaName)
			if err != tt.wantErr {
				t.Errorf("Planisphere.GetSchemaFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planisphere.GetSchemaFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}
