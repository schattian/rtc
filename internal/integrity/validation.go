package integrity

// Validator is any func that returns an error if it doesn't pass the integrity checks
type Validator func() error
