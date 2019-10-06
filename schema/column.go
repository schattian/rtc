package schema

// A Column is the representation of SQL column which defines the structure of the fields that is contains.
type Column struct {
	Name      ColumnName
	Validator func(interface{}) error
}

// ColumnName is the name of a column
type ColumnName string
