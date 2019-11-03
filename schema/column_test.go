package schema

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema/valide"
)

func TestColumn_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name      integrity.ColumnName
		Validator integrity.Validator
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "passes the validation",
			fields:  fields{Validator: valide.String},
			args:    args{val: "anything"},
			wantErr: false,
		},
		{
			name:    "doesnt passes the validation",
			fields:  fields{Validator: valide.Int},
			args:    args{val: "anything"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Column{
				Name:      tt.fields.Name,
				Validator: tt.fields.Validator,
			}
			if err := c.Validate(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Column.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColumn_applyBuiltinValidator(t *testing.T) {
	tests := []struct {
		name          string
		Type          integrity.ValueType
		wantValidator integrity.Validator
		wantType      integrity.ValueType
		wantErr       bool
	}{
		{
			name:          "JSON type",
			Type:          "json",
			wantType:      "json.RawMessage",
			wantValidator: valide.JSON,
			wantErr:       false,
		},
		{
			name:          "BYTES type",
			Type:          "bytes",
			wantType:      "[]byte",
			wantValidator: valide.Bytes,
			wantErr:       false,
		},
		{
			name:          "STR type",
			Type:          "string",
			wantValidator: valide.String,
			wantErr:       false,
		},
		{
			name:          "INT type",
			Type:          "int",
			wantValidator: valide.Int,
			wantErr:       false,
		},
		{
			name:          "FLOAT32 type",
			Type:          "float32",
			wantValidator: valide.Float32,
			wantErr:       false,
		},
		{
			name:          "FLOAT64 type",
			Type:          "float64",
			wantValidator: valide.Float64,
			wantErr:       false,
		},
		{
			name:    "NIL type",
			Type:    "",
			wantErr: true,
		},

		{
			name:    "INVALID type",
			Type:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{Type: tt.Type}
			err := c.applyBuiltinValidator()
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.applyBuiltinValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if c.Validator != nil {
					t.Errorf("Column.applyBuiltinValidator() ASSIGNED a VALIDATOR while it errored: %v", c.Validator)
				}
				if c.Type != tt.Type {
					t.Errorf("Column.applyBuiltinValidator() CHANGED a TYPE while it errored: %v", c.Validator)
				}
				return
			}

			if fmt.Sprintf("%v", c.Validator) != fmt.Sprintf("%v", tt.wantValidator) {
				t.Errorf("Column.applyBuiltinValidator() mismatch VALIDATOR")
			}

			if tt.wantType == "" {
				tt.wantType = tt.Type
			}
			if diff := cmp.Diff(tt.wantType, c.Type); diff != "" {
				t.Errorf("Column.applyBuiltinValidator() mismatch TYPE (-want +got): %s", diff)
			}
		})
	}
}
