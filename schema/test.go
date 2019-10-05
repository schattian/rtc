package schema

import "github.com/sebach1/git-crud/internal/tassist"

var (
	gSchemas goldenSchemas
	gTables  goldenTables
	gColumns goldenColumns
)

func init() {
	tassist.DecodeJsonnet("schemas", &gSchemas)
	tassist.DecodeJsonnet("columns", &gColumns)
	tassist.DecodeJsonnet("tables", &gTables)
}

type goldenSchemas struct {
	Basic     *Schema
	Rare      *Schema
	BasicRare *Schema `json:"basic_rare"`
	Zero      *Schema
}

type goldenTables struct {
	Basic     *Table
	Rare      *Table
	BasicRare *Table `json:"basic_rare"`
	Zero      *Table
}

type goldenColumns struct {
	Basic *Column
	Rare  *Column
	Zero  *Column
}
