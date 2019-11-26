package xerrors

import (
	"errors"
	"reflect"
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
			fields: fields{OriginType: "foo", OriginName: "bar", Err: errBaz},
			want:   "foo validation error: baz. At bar",
		},
		{
			name:   "SKIP ORIGIN information from UNFILLED validationError",
			fields: fields{Err: errBaz},
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

func TestValidationError_Unwrap(t *testing.T) {
	type fields struct {
		OriginType string
		OriginName string
		Err        error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "err is nil",
			fields:  fields{Err: nil},
			wantErr: nil,
		},
		{
			name:    "SKIP ORIGIN information from UNFILLED validationError",
			fields:  fields{Err: errors.New("baz")},
			wantErr: errBaz,
		},
		{
			name:    "UNWRAPs CORRECTly a FULLFILLed validationError",
			fields:  fields{OriginType: "foo", OriginName: "bar", Err: errBaz},
			wantErr: errBaz,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vErr := ValidationError{
				OriginType: tt.fields.OriginType,
				OriginName: tt.fields.OriginName,
				Err:        tt.fields.Err,
			}
			if err := UnwrapValidationError(errors.New(vErr.Error())); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ValidationError.Unwrap() error = %v, wantErr %v.", err, tt.wantErr)
			}
		})
	}
}
