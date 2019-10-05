package schema

// The Schema is the representation of a Database instructive. It uses concepts of SQL.
// The Schema provided by the schema gives the validation structure.
type Schema struct {
	Name      string   `json:"name,omitempty"`
	Blueprint []*Table `json:"blueprint,omitempty"`
}

// ID is an interface made to handle with every type of id types (UUID/GUID/int)
type ID interface{}

// TableName is the name of a table
type TableName string

// ColumnName is the name of a column
type ColumnName string

// The Value is a url-encoded representation of a value
type Value string

// A Table is the representation of SQL table (or Mongo/CQL Collections) which acts as a collection of entities.
type Table struct {
	Name    TableName
	Columns []*Column
}

// A Column is the representation of SQL column which defines the structure of the fields that is contains.
type Column struct {
	Name      ColumnName
	Validator func(*Column) bool
}
