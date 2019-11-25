package integrity

import "testing"

func TestId_IsNil(t *testing.T) {
	tests := []struct {
		name string
		id   Id
		want bool
	}{
		{name: "foo", id: "foo", want: false},
		{name: "nil", id: "", want: true},
		{name: "zero", id: "0", want: true},
		{name: "multiple zeros", id: "0000000000", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsNil(); got != tt.want {
				t.Errorf("Id.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}
