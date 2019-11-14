package schema

import "errors"

var (
	// Schema errs
	errNilBlueprint  = errors.New("the BLUEPRINT is NIL")
	errNilSchema     = errors.New("the SCHEMA cannot be NIL")
	errNilSchemaName = errors.New("the SCHEMA NAME cannot be NIL")

	// Decoding
	errUnallowedExt = errors.New("the EXTension is NOT ALLOWED")

	// Table errs
	errNonexistentTable = errors.New("the TABLE given does NOT EXISTS")
	errForeignTable     = errors.New("the TABLE given does NOT BELONGS to the given SCHEMA")
	errInvalidOptionKey = errors.New("the provided OPTION KEY is INVALID OR NOT DEFINED")
	errNilTableName     = errors.New("the TABLE NAME is NIL")
	errNilColumns       = errors.New("the COLUMNS cannot be NIL")
	errNilTable         = errors.New("the TABLE is NIL")

	// Column errs
	errNonexistentColumn   = errors.New("the COLUMN given does NOT EXISTS")
	errForeignColumn       = errors.New("the COLUMNS given does NOT BELONGS to the given TABLE")
	errNilColumnName       = errors.New("the COLUMN NAME is NIL")
	errNilColumn           = errors.New("the COLUMN is NIL")
	errUnallowedColumnType = errors.New("the COLUMN TYPE is NOT ALLOWED")
	errNilColumnType       = errors.New("the COLUMN TYPE is NIL")

	// Planisphere
	errNotFoundSchema = errors.New("the given SCHEMA NAME is NOT FOUND")
)
