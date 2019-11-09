package git

import (
	"github.com/sebach1/git-crud/integrity"
)

// Options is a bunch of key-value pairs which must must be in accordance to the OptionKeys
// of the table that it belongs to.
// Options is a design error, and in future releases it should be rethinked.
type Options map[integrity.OptionKey]interface{}

// Keys returns the OptionKeys slice
func (opts Options) Keys() (keys []integrity.OptionKey) {
	for k := range opts {
		keys = append(keys, k)
	}
	return
}
