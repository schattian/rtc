package fabric

import (
	"github.com/sebach1/git-crud/internal/test/assist"
	"github.com/sebach1/git-crud/schema"
)

var (
	gSchemas schema.GoldenSchemas
	gTables  schema.GoldenTables
	gColumns schema.GoldenColumns
)

func init() {
	assist.DecodeJsonnet("schemas", &gSchemas)
	assist.DecodeJsonnet("tables", &gTables)
	assist.DecodeJsonnet("columns", &gColumns)
}
