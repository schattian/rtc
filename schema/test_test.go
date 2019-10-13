package schema

import (
	"github.com/sebach1/git-crud/internal/test/assist"
)

var (
	gSchemas goldenSchemas
	gTables  goldenTables
	gColumns goldenColumns
)

func init() {
	assist.DecodeJsonnet("schemas", &gSchemas)
	assist.DecodeJsonnet("columns", &gColumns)
	assist.DecodeJsonnet("tables", &gTables)
}

type goldenSchemas struct {
	Basic     *Schema `json:"basic,omitempty"`
	Rare      *Schema `json:"rare,omitempty"`
	BasicRare *Schema `json:"basic_rare,omitempty"`
	Zero      *Schema `json:"zero,omitempty"`
}

type goldenTables struct {
	Basic     *Table `json:"basic,omitempty"`
	Rare      *Table `json:"rare,omitempty"`
	BasicRare *Table `json:"basic_rare,omitempty"`
	Zero      *Table `json:"zero,omitempty"`
}

type goldenColumns struct {
	Basic *Column `json:"basic,omitempty"`
	Rare  *Column `json:"rare,omitempty"`
	Zero  *Column `json:"zero,omitempty"`
}
