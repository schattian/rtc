package schema

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/integrity"
	"github.com/sebach1/rtc/schema/valide"
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
		name     string
		fields   fields
		args     args
		wantsErr bool // Notice wantSErr (errs are not necessarily std due they're wrapped from the col validator)
	}{
		{
			name:     "passes the validation",
			fields:   fields{Validator: valide.String},
			args:     args{val: "anything"},
			wantsErr: false,
		},
		{
			name:     "doesnt passes the validation",
			fields:   fields{Validator: valide.Int},
			args:     args{val: "anything"},
			wantsErr: true,
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
			if err := c.Validate(tt.args.val); (err != nil) != tt.wantsErr {
				t.Errorf("Column.Validate() error = %v, wantErr %v", err, tt.wantsErr)
			}
		})
	}
}

func TestColumn_applyBuiltinValidator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		Type          integrity.ValueType
		wantValidator integrity.Validator
		wantType      integrity.ValueType
		wantErr       error
	}{
		{
			name:          "JSON type",
			Type:          "json",
			wantType:      "json.RawMessage",
			wantValidator: valide.JSON,
			wantErr:       nil,
		},
		{
			name:          "BYTES type",
			Type:          "bytes",
			wantType:      "[]byte",
			wantValidator: valide.Bytes,
			wantErr:       nil,
		},
		{
			name:          "STR type",
			Type:          "string",
			wantValidator: valide.String,
			wantErr:       nil,
		},
		{
			name:          "INT type",
			Type:          "int",
			wantValidator: valide.Int,
			wantErr:       nil,
		},
		{
			name:          "FLOAT32 type",
			Type:          "float32",
			wantValidator: valide.Float32,
			wantErr:       nil,
		},
		{
			name:          "FLOAT64 type",
			Type:          "float64",
			wantValidator: valide.Float64,
			wantErr:       nil,
		},
		{
			name:    "NIL type",
			Type:    "",
			wantErr: errNilColumnType,
		},

		{
			name:    "INVALId type",
			Type:    "invalid",
			wantErr: errUnallowedColumnType,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Column{Type: tt.Type}
			err := c.applyBuiltinValidator()
			if err != tt.wantErr {
				t.Errorf("Column.applyBuiltinValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if c.Validator != nil {
					t.Errorf("Column.applyBuiltinValidator() ASSIGNED a VALIdATOR while it errored: %v", c.Validator)
				}
				if c.Type != tt.Type {
					t.Errorf("Column.applyBuiltinValidator() CHANGED a TYPE while it errored: %v", c.Validator)
				}
				return
			}

			if fmt.Sprintf("%v", c.Validator) != fmt.Sprintf("%v", tt.wantValidator) {
				t.Errorf("Column.applyBuiltinValidator() mismatch VALIdATOR")
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
