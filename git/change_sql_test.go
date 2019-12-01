package git

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gedex/inflector"
	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/internal/name"
)

func TestChangeSQLColumns(t *testing.T) {
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

	got := chg.SQLColumns()
	sort.Strings(got)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Change.SQLColumns() mismatch (-want +got): %s", diff)
	}
}

func TestChangeSQLTable(t *testing.T) {
	chg := Change{}
	typeOf := reflect.TypeOf(chg)
	want := inflector.Pluralize(name.ToSnakeCase(typeOf.Name()))
	got := chg.SQLTable()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Change.SQLTable() mismatch (-want +got): %s", diff)
	}
}
