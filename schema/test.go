package schema

var ErrSchemaNotFoundInScope = errSchemaNotFoundInScope

// GoldenSchemas is the golden file decoder for the schemas. Exposed to serve to git integration tests
type GoldenSchemas struct {
	Foo    *Schema `json:"foo,omitempty"`
	Bar    *Schema `json:"bar,omitempty"`
	FooBar *Schema `json:"foo_bar,omitempty"`
	Zero   *Schema `json:"zero,omitempty"`
}

// GoldenTables is the golden file decoder for the tables. Exposed to serve to git integration tests
type GoldenTables struct {
	Foo    *Table `json:"foo,omitempty"`
	Bar    *Table `json:"bar,omitempty"`
	FooBar *Table `json:"foo_bar,omitempty"`
	Zero   *Table `json:"zero,omitempty"`
}

// GoldenColumns is the golden file decoder for the columns. Exposed to serve to git integration tests
type GoldenColumns struct {
	Foo  *Column `json:"foo,omitempty"`
	Bar  *Column `json:"bar,omitempty"`
	Zero *Column `json:"zero,omitempty"`
}
