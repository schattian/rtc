package maptypes

import (
	"testing"
)

func TestCommit_GroupBy(t *testing.T) {
	type fields struct {
		ID      int
		Changes []*Change
	}
	type args struct {
		comparator func(*Change, *Change) bool
	}

	alwaysNo := func(*Change, *Change) bool { return false }
	alwaysYes := func(*Change, *Change) bool { return true }
	areSameTable := func(a *Change, b *Change) bool { return a.TableName == b.TableName }

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantQtGrps int
	}{
		{
			name: "all changes are groupable",
			fields: fields{Changes: []*Change{
				&Change{}, &Change{}, &Change{},
			}},
			args:       args{comparator: alwaysYes},
			wantQtGrps: 1,
		},
		{
			name: "all changes are UNgroupable",
			fields: fields{Changes: []*Change{
				&Change{}, &Change{}, &Change{},
			}},
			args:       args{comparator: alwaysNo},
			wantQtGrps: 3,
		},
		{
			name: "changes are groupable if with same tableName",
			fields: fields{Changes: []*Change{
				&Change{TableName: "a"}, &Change{TableName: "a"},
				&Change{TableName: "b"}, &Change{TableName: "b"}, &Change{TableName: "b"},
				&Change{TableName: "c"},
				&Change{},
			}},
			args:       args{comparator: areSameTable},
			wantQtGrps: 4,
		},
		{
			name: "SHUFFLED changes are groupable if with same tableName",
			fields: fields{Changes: []*Change{
				&Change{},
				&Change{TableName: "b"}, &Change{TableName: "c"}, &Change{TableName: "a"},
				&Change{TableName: "a"}, &Change{TableName: "a"}, &Change{TableName: "a"},
				&Change{TableName: "b"}, &Change{TableName: "b"}, &Change{TableName: "a"},
			}},
			args:       args{comparator: areSameTable},
			wantQtGrps: 4,
		},
		{
			name: "changes are groupable if are compatible",
			fields: fields{Changes: []*Change{
				&Change{TableName: "a"},
				&Change{TableName: "a"},
				&Change{TableName: "a", EntityID: 1}, &Change{TableName: "a", EntityID: 1},
				&Change{TableName: "b", EntityID: 5},
				&Change{TableName: "b", EntityID: 3},
			}},
			args:       args{comparator: IsCompatibleWith},
			wantQtGrps: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comm := &Commit{
				ID:      tt.fields.ID,
				Changes: tt.fields.Changes,
			}

			if gotQtGrps := len(comm.GroupBy(tt.args.comparator)); gotQtGrps != tt.wantQtGrps {
				t.Errorf("Commit.GroupBy() = %v, want %v", gotQtGrps, tt.wantQtGrps)
			}
		})
	}
}

func Test_checkIntInSlice(t *testing.T) {
	type args struct {
		slice []int
		elem  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "is in slice", args: args{slice: []int{1, 2}, elem: 2}, want: true},
		{name: "is NOT in slice", args: args{slice: []int{1, 2}, elem: 3}, want: false},
		{name: "slice is empty", args: args{slice: []int{}, elem: 3}, want: false},
		{name: "elem is zero-value", args: args{slice: []int{}, elem: 0}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIntInSlice(tt.args.slice, tt.args.elem); got != tt.want {
				t.Errorf("checkIntInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
