package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gedex/inflector"
	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/name"
)

func TestBranch_columns(t *testing.T) {
	b := Branch{}
	exclusions := []string{"Index", "Credentials"}
	typeOf := reflect.TypeOf(b)
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

	got := b.columns()
	sort.Strings(got)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Branch.columns() mismatch (-want +got): %s", diff)
	}
}

func TestBranch_table(t *testing.T) {
	b := Branch{}
	typeOf := reflect.TypeOf(b)
	want := inflector.Pluralize(name.ToSnakeCase(typeOf.Name()))
	got := b.table()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Branch.table() mismatch (-want +got): %s", diff)
	}
}
