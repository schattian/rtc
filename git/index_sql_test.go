package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gedex/inflector"
	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/internal/name"
)

func TestIndexSQLColumns(t *testing.T) {
	idx := Index{}
	exclusions := []string{"Changes"}
	typeOf := reflect.TypeOf(idx)
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

	got := idx.SQLColumns()
	sort.Strings(got)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Index.SQLColumns() mismatch (-want +got): %s", diff)
	}
}

func TestIndexSQLTable(t *testing.T) {
	idx := Index{}
	typeOf := reflect.TypeOf(idx)
	want := inflector.Pluralize(name.ToSnakeCase(typeOf.Name()))
	got := idx.SQLTable()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Index.SQLTable() mismatch (-want +got): %s", diff)
	}
}
