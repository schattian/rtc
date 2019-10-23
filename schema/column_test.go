package schema

import (
	"testing"

	"github.com/sebach1/git-crud/schema/valide"

	"github.com/sebach1/git-crud/integrity"
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
