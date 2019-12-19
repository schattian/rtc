package schema

import (
	"testing"

	"github.com/mitchellh/copystructure"
)

func (sch *Schema) copy(t *testing.T) *Schema {
	t.Helper()
	new, err := copystructure.Copy(sch)
	if err != nil {
		t.Fatalf("could not be able to copy schema: %v", err)
	}
	return new.(*Schema)
}
