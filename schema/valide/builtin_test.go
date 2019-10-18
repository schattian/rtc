package valide

import (
	"testing"
)

func TestString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "string",
			val:     "test",
			wantErr: false,
		},
		{
			name:    "int",
			val:     3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := new(String).Validate(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "int",
			val:     3,
			wantErr: false,
		},
		{
			name:    "float32",
			val:     float32(3),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := new(Int).Validate(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("Int() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "float32",
			val:     float32(3),
			wantErr: false,
		},
		{
			name:    "float64",
			val:     float64(3),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := new(Float32).Validate(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("Float32() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "float64",
			val:     float64(3),
			wantErr: false,
		},

		{
			name:    "float32",
			val:     float32(3),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := new(Float64).Validate(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("Float64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "json",
			val:     []byte(`{"id":"1"}`),
			wantErr: false,
		},
		{
			name:    "int",
			val:     3,
			wantErr: true,
		},
		{
			name:    "urlparam",
			val:     []byte(`id=1`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := new(JSON).Validate(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("JSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "json",
			val:     []byte(`{"id":"1"}`),
			wantErr: false,
		},
		{
			name:    "urlparam",
			val:     []byte(`id=1`),
			wantErr: false,
		},
		{
			name:    "int",
			val:     3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := new(Bytes).Validate(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("Bytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
