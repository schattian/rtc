package xerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")
	errBaz = errors.New("baz")
)

func TestNewMultiErr(t *testing.T) {
	tests := []struct {
		name     string
		errs     []error
		wantMErr MultiErr
	}{
		{
			name:     "multiple normal errs",
			errs:     []error{errFoo, errBar, errBaz},
			wantMErr: MultiErr{errFoo, errBar, errBaz},
		},
		{
			name:     "nil errs",
			errs:     []error{nil, nil, nil},
			wantMErr: MultiErr{},
		},
		{
			name:     "nil and normal errs",
			errs:     []error{nil, errFoo, nil},
			wantMErr: MultiErr{errFoo},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMErrs := NewMultiErr(tt.errs...)
			if diff := cmp.Diff(tt.wantMErr.Error(), gotMErrs.Error()); diff != "" {
				t.Errorf("NewMultiErr() = mismatch(-want +got): %s", diff)
			}
		})
	}
}

func TestMultiErr_Error(t *testing.T) {
	tests := []struct {
		name string
		errs MultiErr
		want string
	}{
		{
			name: "multiple normal errs",
			errs: []error{errFoo, errBar, errBaz},
			want: errFoo.Error() + ErrorsSeparator + errBar.Error() + ErrorsSeparator + errBaz.Error() + ErrorsSeparator,
		},
		{
			name: "nil and normal errs",
			errs: []error{nil, errFoo, nil},
			want: fmt.Sprintf("%s%s", errFoo, ErrorsSeparator),
		},
		{
			name: "nil  errs",
			errs: []error{nil, nil, nil},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.errs.Error()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MultiErr.Error() =  mismatch (-want +got): %s", diff)
			}
		})
	}
}
