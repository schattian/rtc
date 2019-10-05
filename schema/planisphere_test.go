package schema

import (
	"reflect"
	"testing"
)

func Test_preciseTableErr(t *testing.T) {
	type args struct {
		tableName TableName
	}

	tests := []struct {
		name string
		args args
		psph Planisphere
		want error
	}{
		{
			name: "given tableName is in a schema",
			args: args{tableName: gTables.Basic.Name},
			psph: Planisphere{gSchemas.Basic},
			want: errForeignTable,
		},
		{
			name: "given tableName doenst exists on any scoped schema",
			args: args{gTables.Zero.Name},
			psph: Planisphere{gSchemas.Basic, gSchemas.Rare},
			want: errUnexistantTable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.psph.preciseTableErr(tt.args.tableName); err != tt.want {
				t.Errorf("preciseTableErr() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestPlanisphere_GetSchemaFromName(t *testing.T) {
	type args struct {
		schemaName string
	}

	tests := []struct {
		name    string
		psph    Planisphere
		args    args
		want    *Schema
		wantErr bool
	}{
		{
			name:    "blank schemaName and having blank schemas",
			args:    args{gSchemas.Zero.Name},
			psph:    Planisphere{gSchemas.Zero, gSchemas.Zero},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "giving a single schema matching the given name",
			args:    args{gSchemas.Basic.Name},
			psph:    Planisphere{gSchemas.Basic},
			want:    gSchemas.Basic,
			wantErr: false,
		},
		{
			name:    "giving multiple gSchemas matching the given name it matches the first",
			args:    args{gSchemas.Basic.Name},
			psph:    Planisphere{gSchemas.Basic, gSchemas.BasicRare},
			want:    gSchemas.Basic,
			wantErr: false,
		},
		{
			name:    "giving no schema",
			args:    args{gSchemas.Basic.Name},
			psph:    Planisphere{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.psph.GetSchemaFromName(tt.args.schemaName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Planisphere.GetSchemaFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planisphere.GetSchemaFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}
