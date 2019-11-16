package integrity

import (
	"testing"
)

func TestCRUD_ToHTTPVerb(t *testing.T) {
	t.Parallel()
	tests := []struct {
		crud CRUD
		want string
	}{
		{crud: "create", want: "POST"},
		{crud: "update", want: "PUT"},
		{crud: "retrieve", want: "GET"},
		{crud: "delete", want: "DELETE"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.crud), func(t *testing.T) {
			t.Parallel()
			if got := tt.crud.ToHTTPVerb(); tt.want != got {
				t.Errorf("CRUD.ToHTTPVerb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCRUD_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		crud CRUD
		want error
	}{
		{crud: "create", want: nil},
		{crud: "update", want: nil},
		{crud: "retrieve", want: nil},
		{crud: "delete", want: nil},
		{crud: "foo", want: errInvalidCRUD},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.crud), func(t *testing.T) {
			t.Parallel()
			if got := tt.crud.Validate(); tt.want != got {
				t.Errorf("CRUD.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
