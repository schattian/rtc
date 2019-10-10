package schema

import (
	"errors"
	"testing"

	"github.com/sebach1/git-crud/internal/integrity"
)

func TestColumn_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name      integrity.ColumnName
		Validator func(interface{}) error
	}
	type args struct {
		val interface{}
	}
	neverErr := func(val interface{}) error { return nil }
	alwaysErr := func(val interface{}) error { return errors.New("err") }
	validateStr := func(val interface{}) error {
		if _, ok := val.(string); ok {
			return nil
		}
		return errors.New("err")
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "never errs",
			fields:  fields{Validator: neverErr},
			args:    args{val: "anything"},
			wantErr: false,
		},
		{
			name:    "str validator",
			fields:  fields{Validator: validateStr},
			args:    args{val: "anything"},
			wantErr: false,
		},
		{
			name:    "always errs",
			fields:  fields{Validator: alwaysErr},
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
