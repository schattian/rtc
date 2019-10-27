package fabric

import "errors"

var errUnnamedCol = errors.New("the SCHEMA CONTAINS an UNNAMED COLUMN")
var errFmtUntypedColumn = "the TYPE of COLUMN %v cannot be NIL"

var errNilSchema = errors.New("the SCHEMA cant be NIL")
var errUnnamedSchema = errors.New("the SCHEMA NAME cant be NIL")
