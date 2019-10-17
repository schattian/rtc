package schema

// GoldenSchemas is the golden file decoder for the schemas. Exposed to serve to git integration tests
type GoldenSchemas struct {
	Basic     *Schema `json:"basic,omitempty"`
	Rare      *Schema `json:"rare,omitempty"`
	BasicRare *Schema `json:"basic_rare,omitempty"`
	Zero      *Schema `json:"zero,omitempty"`
}

// GoldenTables is the golden file decoder for the tables. Exposed to serve to git integration tests
type GoldenTables struct {
	Basic     *Table `json:"basic,omitempty"`
	Rare      *Table `json:"rare,omitempty"`
	BasicRare *Table `json:"basic_rare,omitempty"`
	Zero      *Table `json:"zero,omitempty"`
}

// GoldenColumns is the golden file decoder for the columns. Exposed to serve to git integration tests
type GoldenColumns struct {
	Basic *Column `json:"basic,omitempty"`
	Rare  *Column `json:"rare,omitempty"`
	Zero  *Column `json:"zero,omitempty"`
}
