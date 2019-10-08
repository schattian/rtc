package git

import "errors"

var (
	// Changes
	errDuplicatedChg     = errors.New("the change is ALREADY COMMITED")
	errUnclassifiableChg = errors.New("the change DOESNT respect any PATTERN and thus CANNOT be CLASSIFIABLE")

	errNilTable = errors.New("change's TABLE cannot be NIL")

	errNotNilColumn = errors.New("change's COLUMN cannot be NOT NIL")
	errNilColumn    = errors.New("change's COLUMN cannot be NIL")

	errNilEntityID    = errors.New("the ENTITY_ID is NIL")
	errNotNilEntityID = errors.New("the ENTITY_ID is NOT NIL")

	errNilValue    = errors.New("the VALUE cannot be NIL")
	errNotNilValue = errors.New("the VALUE cannot be NOT NIL")

	// Commit
	errMixedTypes  = errors.New("the TYPES over the commit are MIXED")
	errMixedTables = errors.New("the TABLES over the commit are MIXED")

	// Community
	errNotFoundSchema = errors.New("the SCHEMA NAME provided is NOT FOUND")
	errNilCommunity   = errors.New("the COMMUNITY cannot be NIL")
)
