package git

import (
	"github.com/sebach1/git-crud/internal/integrity"
)

// Options is a bunch of key-value pairs which must must be in accordance to the OptionKeys
// of the table that it belongs to.
// Options is a design error, and in future releases it should be rethinked.
type Options map[integrity.OptionKey]interface{}
