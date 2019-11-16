package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gedex/inflector"
	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/name"
)

func TestChange_columns(t *testing.T) {
	chg := Change{}
	exclusions := []string{}
	typeOf := reflect.TypeOf(chg)
	var want []string
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if isExcluded(exclusions, field.Name) {
			continue
		}
		col := name.ToSnakeCase(field.Name)
		want = append(want, col)
	}
	sort.Strings(want)

	got := chg.columns()
	sort.Strings(got)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Change.columns() mismatch (-want +got): %s", diff)
	}
}

func TestChange_table(t *testing.T) {
	chg := Change{}
	typeOf := reflect.TypeOf(chg)
	want := inflector.Pluralize(name.ToSnakeCase(typeOf.Name()))
	got := chg.table()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Change.table() mismatch (-want +got): %s", diff)
	}
}
