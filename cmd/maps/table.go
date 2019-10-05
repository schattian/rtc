package maps

// A Table is the representation of SQL table (or Mongo/CQL Collections) which acts as a collection of entities.
// Note: every name/columnName should be in snake_case (to be used as a std url param)
// Note 2: every column should be in snake_case (to be used as a std url param)
type Table struct {
	Name    string
	Columns []string
}
