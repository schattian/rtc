package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/name"
)

func TestChange_Columns(t *testing.T) {
	chg := Change{}
	exclusions := []string{}
	typeOf := reflect.TypeOf(chg)
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

	got := chg.Columns()
	sort.Strings(got)
	if diff := cmp.Diff(got, cols); diff != "" {
		t.Errorf("Change.Columns() mismatch (-want +got): %s", diff)
	}
}
