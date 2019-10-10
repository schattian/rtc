package valide

import "errors"

// String tries to convert an interface to a string value; if cant it returns an err
func String(v interface{}) error {
	if _, ok := v.(string); !ok {
		return errors.New("the value isn't a valid string")
	}
	return nil
}
