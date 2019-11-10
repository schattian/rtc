package git

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/integrity"
)

func TestCommit_GroupBy(t *testing.T) {
	t.Parallel()
	type fields struct {
		ID      int64
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
				gChanges.Zero, gChanges.Basic.None, gChanges.Basic.EntityID,
			}},
			args:       args{comparator: alwaysYes},
			wantQtGrps: 1,
		},
		{
			name: "all changes are UNgroupable",
			fields: fields{Changes: []*Change{
				gChanges.Zero, gChanges.Basic.None, gChanges.Basic.EntityID,
			}},
			args:       args{comparator: alwaysNo},
			wantQtGrps: 3,
		},
		{
			name: "changes are groupable if with same tableName",
			fields: fields{Changes: []*Change{
				gChanges.Basic.None, gChanges.Basic.Create, gChanges.Basic.None, gChanges.Rare.TableName,
				gChanges.Basic.TableName,
			}},
			args:       args{comparator: areSameTable},
			wantQtGrps: 2,
		},
		{
			name: "changes are groupable if are compatible",
			fields: fields{Changes: []*Change{
				gChanges.Basic.None, gChanges.Basic.ColumnName,
				gChanges.Rare.None,
				gChanges.Zero,
				gChanges.Basic.TableName,
			}},
			args:       args{comparator: AreCompatible},
			wantQtGrps: 4,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := checkIntInSlice(tt.args.slice, tt.args.elem); got != tt.want {
				t.Errorf("checkIntInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_ToMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		comm *Commit
		want map[string]interface{}
	}{
		{
			name: "CREATE commit with multiple columns",
			comm: &Commit{Changes: []*Change{gChanges.Basic.Create}},
			want: map[string]interface{}{
				string(gChanges.Basic.Create.ColumnName): gChanges.Basic.Create.StrValue,
			},
		},
		{
			name: "RETRIEVE commit",
			comm: &Commit{Changes: []*Change{gChanges.Basic.Delete}},
			want: map[string]interface{}{
				"id": gChanges.Basic.None.EntityID,
			},
		},
		{
			name: "UPDATE commit with multiple column changes",
			comm: &Commit{Changes: []*Change{gChanges.Basic.None, gChanges.Basic.ColumnName}},
			want: map[string]interface{}{
				"id":                                   gChanges.Basic.None.EntityID,
				string(gChanges.Basic.None.ColumnName): gChanges.Basic.None.StrValue,
				string(gChanges.Basic.ColumnName.ColumnName): gChanges.Basic.ColumnName.StrValue,
			},
		},
		{
			name: "DELETE commit",
			comm: &Commit{Changes: []*Change{gChanges.Basic.Delete}},
			want: map[string]interface{}{
				"id": gChanges.Basic.None.EntityID,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if diff := cmp.Diff(tt.want, tt.comm.ToMap()); diff != "" {
				t.Errorf("Commmit.ToMap() mismatch (-want +got): %s", diff)
			}

		})
	}
}

func TestCommit_TableName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		comm          *Commit
		wantTableName integrity.TableName
		wantErr       bool
	}{
		{
			name: "changes contains the same single table",
			comm: &Commit{Changes: []*Change{
				gChanges.Basic.None, gChanges.Rare.TableName,
			}},
			wantTableName: gChanges.Basic.None.TableName,
			wantErr:       false,
		},
		{
			name: "changes contains mixed tables",
			comm: &Commit{Changes: []*Change{
				gChanges.Basic.None, gChanges.Rare.None,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotTableName, err := tt.comm.TableName()
			if (err != nil) != tt.wantErr {
				t.Errorf("Commit.TableName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTableName, tt.wantTableName) {
				t.Errorf("Commit.TableName() = %v, want %v", gotTableName, tt.wantTableName)
			}
		})
	}
}

func TestCommit_Type(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		comm         *Commit
		wantCommType integrity.CRUD
		wantErr      bool
	}{
		{
			name: "changes contains the same single table",
			comm: &Commit{Changes: []*Change{
				gChanges.Basic.Create, gChanges.Basic.Create,
			}},
			wantCommType: gChanges.Basic.Create.Type,
			wantErr:      false,
		},
		{
			name: "changes contains mixed types",
			comm: &Commit{Changes: []*Change{
				gChanges.Basic.Delete, gChanges.Basic.Create,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotCommType, err := tt.comm.Type()
			if (err != nil) != tt.wantErr {
				t.Errorf("Commit.Type() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCommType, tt.wantCommType) {
				t.Errorf("Commit.Type() = %v, want %v", gotCommType, tt.wantCommType)
			}
		})
	}
}
