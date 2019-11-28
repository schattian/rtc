package schema

import (
	"github.com/sebach1/rtc/internal/test/assist"
)

var (
	gSchemas GoldenSchemas
	gTables  GoldenTables
	gColumns GoldenColumns
)

func init() {
	assist.DecodeJsonnet("schemas", &gSchemas)
	assist.DecodeJsonnet("columns", &gColumns)
	assist.DecodeJsonnet("tables", &gTables)
}
