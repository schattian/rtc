package git

import "errors"

var (
	// Changes
	errInvalidChangeId   = errors.New("the change Id is NOT AN Id TYPE")
	errDuplicatedChg     = errors.New("the change is ALREADY COMMITTED")
	errUnclassifiableChg = errors.New("the change DOESN'T respect any PATTERN and thus CANNOT be CLASSIFIABLE")
	errUnsafeValueType   = errors.New("the given value cannot be safety typed")
	errNilOptionKey      = errors.New("the given OPTION KEY is NIL")

	// Table
	errNilTable = errors.New("change's TABLE cannot be NIL")

	// Column
	errNilColumn    = errors.New("change's COLUMN cannot be NIL")
	errNotNilColumn = errors.New("change's COLUMN cannot be NOT NIL")

	// EntityId
	errNilEntityId    = errors.New("the ENTITY_Id is NIL")
	errNotNilEntityId = errors.New("the ENTITY_Id is NOT NIL")

	// Value
	errNilValue    = errors.New("the VALUE cannot be NIL")
	errNotNilValue = errors.New("the VALUE cannot be NOT NIL")

	// Commit
	errInvalidCommitId = errors.New("the commit Id is NOT AN Id TYPE")
	errMixedTypes      = errors.New("the TYPES over the commit are MIXED")
	errMixedTables     = errors.New("the TABLES over the commit are MIXED")
	errMixedOpts       = errors.New("the OPTIONS over the commit are MIXED")

	// Community
	errSchemaNotFoundInCommunity = errors.New("the SCHEMA NAME provided is NOT FOUND in the community")
	errNilCommunity              = errors.New("the COMMUNITY cannot be NIL")

	// Team
	errTableInUse      = errors.New("the TABLE is ALREADY IN USE by a member")
	errNoCollaborators = errors.New("there are NOT COLLABORATORS to achieve this TABLE")
	errNoMembers       = errors.New("there are NOT MEMBERS to achieve this TABLE")

	// Owner
	errErroredMerge = errors.New("the MERGE was ERRORED")
	errNilProject   = errors.New("the PROJECT is NIL")
	errEmptyProject = errors.New("the PROJECT does NOT contain ANY SCHEMA")
)
