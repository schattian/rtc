package valide

import (
	"encoding/json"
	"errors"

	"github.com/sebach1/git-crud/integrity"
)

// String validates string
type String struct{}

// NativeType returns the stringified native type of any value that can pass the validation
func (s *String) NativeType() integrity.ValueType {
	return "string"
}

// Validate tries to convert an interface to a string value; if cant it returns an err
func (s *String) Validate(v interface{}) error {
	if _, ok := v.(string); !ok {
		return errors.New("the value isn't a valid string")
	}
	return nil
}

// Int validates int
type Int struct{}

// NativeType returns the stringified native type of any value that can pass the validation
func (i *Int) NativeType() integrity.ValueType {
	return "int"
}

// Validate tries to convert an interface to a int value; if cant it returns an err
func (i *Int) Validate(v interface{}) error {
	if _, ok := v.(int); !ok {
		return errors.New("the value isn't a valid int")
	}
	return nil
}

// Float32 validates float32
type Float32 struct{}

// NativeType returns the stringified native type of any value that can pass the validation
func (f32 *Float32) NativeType() integrity.ValueType {
	return "float32"
}

// Validate tries to convert an interface to a float32 value; if cant it returns an err
func (f32 *Float32) Validate(v interface{}) error {
	if _, ok := v.(float32); !ok {
		return errors.New("the value isn't a valid float32")
	}
	return nil
}

// Float64 validates float64
type Float64 struct{}

// NativeType returns the stringified native type of any value that can pass the validation
func (f64 *Float64) NativeType() integrity.ValueType {
	return "float64"
}

// Validate tries to convert an interface to a float64 value; if cant it returns an err
func (f64 *Float64) Validate(v interface{}) error {
	if _, ok := v.(float64); !ok {
		return errors.New("the value isn't a valid float64")
	}
	return nil
}

// JSON validates json.RawMessage
type JSON struct{}

// NativeType returns the stringified native type of any value that can pass the validation
func (js *JSON) NativeType() integrity.ValueType {
	return "json"
}

// Validate tries to unmarshal an interface to a json value; if cant it returns an err
func (js *JSON) Validate(v interface{}) error {
	byVal, ok := v.([]byte)
	if !ok {
		return errors.New("the value isn't a valid []byte")
	}

	var validator struct{}
	err := json.Unmarshal(byVal, &validator)
	if err != nil {
		return errors.New("the value isn't a valid JSON")
	}
	return nil
}

// Bytes validates []byte
type Bytes struct{}

// NativeType returns the stringified native type of any value that can pass the validation
func (by *Bytes) NativeType() integrity.ValueType {
	return "bytes"
}

// Validate tries to convert an interface to a []bytes value; if cant it returns an err
func (by *Bytes) Validate(v interface{}) error {
	if _, ok := v.([]byte); !ok {
		return errors.New("the value isn't a valid []byte")
	}
	return nil
}
