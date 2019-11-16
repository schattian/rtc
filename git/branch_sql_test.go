package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/name"
)

func TestBranch_Columns(t *testing.T) {
	b := Branch{}
	exclusions := []string{"Index", "Credentials"}
	typeOf := reflect.TypeOf(b)
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

	got := b.Columns()
	sort.Strings(got)
	if diff := cmp.Diff(got, cols); diff != "" {
		t.Errorf("Branch.Columns() mismatch (-want +got): %s", diff)
	}
}
