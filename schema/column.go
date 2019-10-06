package schema

// A Column is the representation of SQL column which defines the structure of the fields that is contains.
type Column struct {
	Name      ColumnName
	Validator func(interface{}) error
}

// ColumnName is the name of a column
type ColumnName string

// Validate wraps the column validator func and returns its9 result
func (c *Column) Validate(val interface{}) error {
	err := c.Validator(val)
	if err != nil {
		return err
	}
	return nil
}
