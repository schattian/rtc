package integrity

// Validator is any func that returns an error if it doesn't pass the integrity checks for a given interface
type Validator interface {
	NativeType() string
	Validate(v interface{}) error
}
