package integrity

// Validator is a func that returns an error if it doesn't pass the integrity checks for a given interface
type Validator func(interface{}) error
