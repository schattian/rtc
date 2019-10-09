package git

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/integrity"
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
				gChanges.Zero, gChanges.Regular.None, gChanges.Regular.EntityID,
			}},
			args:       args{comparator: alwaysYes},
			wantQtGrps: 1,
		},
		{
			name: "all changes are UNgroupable",
			fields: fields{Changes: []*Change{
				gChanges.Zero, gChanges.Regular.None, gChanges.Regular.EntityID,
			}},
			args:       args{comparator: alwaysNo},
			wantQtGrps: 3,
		},
		{
			name: "changes are groupable if with same tableName",
			fields: fields{Changes: []*Change{
				gChanges.Regular.None, gChanges.Regular.Create, gChanges.Regular.None, gChanges.Rare.TableName,
				gChanges.Regular.TableName,
			}},
			args:       args{comparator: areSameTable},
			wantQtGrps: 2,
		},
		{
			name: "changes are groupable if are compatible",
			fields: fields{Changes: []*Change{
				gChanges.Regular.None, gChanges.Regular.ColumnName,
				gChanges.Rare.None,
				gChanges.Zero,
				gChanges.Regular.TableName,
			}},
			args:       args{comparator: AreCompatible},
			wantQtGrps: 4,
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

func TestCommit_Add(t *testing.T) {
	type args struct {
		chg *Change
	}
	tests := []struct {
		name    string
		comm    *Commit
		args    args
		want    error
		newComm *Commit
	}{
		{
			name: "change was already added",
			comm: &Commit{Changes: []*Change{gChanges.Regular.None}},
			args: args{chg: gChanges.Regular.None},
			want: errDuplicatedChg,
		},
		{
			name:    "both identical untracked commits",
			comm:    &Commit{Changes: []*Change{gChanges.Regular.Create}},
			args:    args{chg: gChanges.Regular.Create},
			newComm: &Commit{Changes: []*Change{gChanges.Regular.Create, gChanges.Regular.Create}},
		},
		{
			name: "table inconsistency",
			comm: &Commit{Changes: []*Change{gChanges.Regular.None}},
			args: args{chg: gChanges.Inconsistent.TableName},
			want: errNilTable,
		},
		{
			name:    "change modifies the value of a change already in the commit",
			comm:    &Commit{Changes: []*Change{gChanges.Regular.None}},
			args:    args{chg: gChanges.Regular.StrValue},
			newComm: &Commit{Changes: []*Change{gChanges.Regular.StrValue}},
		},
		{
			name:    "change modifies different col of same schema",
			comm:    &Commit{Changes: []*Change{gChanges.Regular.None}},
			args:    args{chg: gChanges.Regular.ColumnName},
			newComm: &Commit{Changes: []*Change{gChanges.Regular.None, gChanges.Regular.ColumnName}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldComm := tt.comm
			if tt.args.chg == nil {
				panic(tt.name)
			}
			err := tt.comm.Add(tt.args.chg)
			if err != tt.want {
				t.Errorf("Commit.Add() error = %v, wantErr %v", err, tt.want)
			}
			if err != nil {
				if diff := cmp.Diff(oldComm, tt.comm); diff != "" {
					t.Errorf("Commit.Add() errored mismatch (-want +got): %s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.newComm, tt.comm); diff != "" {
				t.Errorf("Commit.Add() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestCommit_Rm(t *testing.T) {
	type args struct {
		chg *Change
	}
	tests := []struct {
		name    string
		comm    *Commit
		args    args
		wantErr bool
	}{
		{
			name:    "given change doesnt belongs to the commit",
			comm:    &Commit{Changes: []*Change{gChanges.Regular.None}},
			args:    args{chg: gChanges.Rare.None},
			wantErr: true,
		},
		{
			name:    "successfully remove",
			comm:    &Commit{Changes: []*Change{gChanges.Regular.None}},
			args:    args{chg: gChanges.Regular.None},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.comm.Rm(tt.args.chg); (err != nil) != tt.wantErr {
				t.Errorf("Commit.Rm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommit_ToMap(t *testing.T) {
	tests := []struct {
		name string
		comm *Commit
		want map[string]interface{}
	}{
		{
			name: "CREATE commit with multiple columns",
			comm: &Commit{Changes: []*Change{gChanges.Regular.Create}},
			want: map[string]interface{}{
				string(gChanges.Regular.Create.ColumnName): gChanges.Regular.Create.StrValue,
			},
		},
		{
			name: "RETRIEVE commit",
			comm: &Commit{Changes: []*Change{gChanges.Regular.Delete}},
			want: map[string]interface{}{
				"id": gChanges.Regular.None.EntityID,
			},
		},
		{
			name: "UPDATE commit with multiple column changes",
			comm: &Commit{Changes: []*Change{gChanges.Regular.None, gChanges.Regular.ColumnName}},
			want: map[string]interface{}{
				"id":                                     gChanges.Regular.None.EntityID,
				string(gChanges.Regular.None.ColumnName): gChanges.Regular.None.StrValue,
				string(gChanges.Regular.ColumnName.ColumnName): gChanges.Regular.ColumnName.StrValue,
			},
		},
		{
			name: "DELETE commit",
			comm: &Commit{Changes: []*Change{gChanges.Regular.Delete}},
			want: map[string]interface{}{
				"id": gChanges.Regular.None.EntityID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.comm.ToMap()); diff != "" {
				t.Errorf("Commmit.ToMap() mismatch (-want +got): %s", diff)
			}

		})
	}
}

func TestCommit_TableName(t *testing.T) {

	tests := []struct {
		name          string
		comm          *Commit
		wantTableName integrity.TableName
		wantErr       bool
	}{
		{
			name: "changes contains the same single table",
			comm: &Commit{Changes: []*Change{
				gChanges.Regular.None, gChanges.Rare.TableName,
			}},
			wantTableName: gChanges.Regular.None.TableName,
			wantErr:       false,
		},
		{
			name: "changes contains mixed tables",
			comm: &Commit{Changes: []*Change{
				gChanges.Regular.None, gChanges.Rare.None,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

	tests := []struct {
		name         string
		comm         *Commit
		wantCommType integrity.CRUD
		wantErr      bool
	}{
		{
			name: "changes contains the same single table",
			comm: &Commit{Changes: []*Change{
				gChanges.Regular.Create, gChanges.Regular.Create,
			}},
			wantCommType: gChanges.Regular.Create.Type,
			wantErr:      false,
		},
		{
			name: "changes contains mixed types",
			comm: &Commit{Changes: []*Change{
				gChanges.Regular.Delete, gChanges.Regular.Create,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
