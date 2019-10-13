package schema

import "github.com/sebach1/git-crud/internal/integrity"

// A Column is the representation of SQL column which defines the structure of the fields that is contains.
type Column struct {
	Name      integrity.ColumnName
	Validator integrity.Validator
}

// Validate wraps the column validator func and returns its result
func (c *Column) Validate(val interface{}) error {
	if c.Validator == nil {
		return nil
	}
	err := c.Validator(val)
	if err != nil {
		return err
	}
	return nil
}
