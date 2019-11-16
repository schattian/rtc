package integrity

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	type fields struct {
		OriginType string
		OriginName string
		Err        error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "WRAPs CORRECTly a FULLFILLed validationError",
			fields: fields{OriginType: "foo", OriginName: "bar", Err: errors.New("baz")},
			want:   "foo validation error: baz. At bar",
		},
		{
			name:   "SKIP ORIGIN information from UNFILLED validationError",
			fields: fields{Err: errors.New("baz")},
			want:   "validation error: baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vErr := &ValidationError{
				OriginType: tt.fields.OriginType,
				OriginName: tt.fields.OriginName,
				Err:        tt.fields.Err,
			}
			if got := vErr.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
