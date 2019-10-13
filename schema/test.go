package schema

type GoldenSchemas struct {
	Basic     *Schema `json:"basic,omitempty"`
	Rare      *Schema `json:"rare,omitempty"`
	BasicRare *Schema `json:"basic_rare,omitempty"`
	Zero      *Schema `json:"zero,omitempty"`
}

type GoldenTables struct {
	Basic     *Table `json:"basic,omitempty"`
	Rare      *Table `json:"rare,omitempty"`
	BasicRare *Table `json:"basic_rare,omitempty"`
	Zero      *Table `json:"zero,omitempty"`
}

type GoldenColumns struct {
	Basic *Column `json:"basic,omitempty"`
	Rare  *Column `json:"rare,omitempty"`
	Zero  *Column `json:"zero,omitempty"`
}
