package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/name"
)

func TestCommit_Columns(t *testing.T) {
	comm := Commit{}
	exclusions := []string{"Reviewer", "Changes"}
	typeOf := reflect.TypeOf(comm)
	var cols []string
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if isExcluded(exclusions, field.Name) {
			continue
		}
		col := name.ToSnakeCase(field.Name)
		cols = append(cols, col)
	}
	sort.Strings(cols)

	got := comm.Columns()
	sort.Strings(got)
	if diff := cmp.Diff(got, cols); diff != "" {
		t.Errorf("Commit.Columns() mismatch (-want +got): %s", diff)
	}
}
