package valide

import (
	"encoding/json"
	"errors"
)

// String tries to convert an interface to a string value; if cant it returns an err
func String(v interface{}) error {
	if _, ok := v.(string); !ok {
		return errors.New("the value isn't a valid string")
	}
	return nil
}

// Int tries to convert an interface to a int value; if cant it returns an err
func Int(v interface{}) error {
	if _, ok := v.(int); !ok {
		return errors.New("the value isn't a valid int")
	}
	return nil
}

// Float32 tries to convert an interface to a float32 value; if cant it returns an err
func Float32(v interface{}) error {
	if _, ok := v.(float32); !ok {
		return errors.New("the value isn't a valid float32")
	}
	return nil
}

// Float64 tries to convert an interface to a float64 value; if cant it returns an err
func Float64(v interface{}) error {
	if _, ok := v.(float64); !ok {
		return errors.New("the value isn't a valid float64")
	}
	return nil
}

// JSON tries to unmarshal an interface to a json value; if cant it returns an err
func JSON(v interface{}) error {
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

// Bytes tries to convert an interface to a []bytes value; if cant it returns an err
func Bytes(v interface{}) error {
	if _, ok := v.([]byte); !ok {
		return errors.New("the value isn't a valid []byte")
	}
	return nil
}
