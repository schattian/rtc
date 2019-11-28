package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gedex/inflector"
	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/internal/name"
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

	got := b.SQLColumns()
	sort.Strings(got)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Branch.SQLColumns() mismatch (-want +got): %s", diff)
	}
}

func TestBranch_table(t *testing.T) {
	b := Branch{}
	typeOf := reflect.TypeOf(b)
	want := inflector.Pluralize(name.ToSnakeCase(typeOf.Name()))
	got := b.SQLTable()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Branch.SQLTable() mismatch (-want +got): %s", diff)
	}
}
