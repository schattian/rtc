package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIdx_Add(t *testing.T) {
	t.Parallel()
	type args struct {
		chg *Change
	}
	tests := []struct {
		name    string
		idx     *Index
		args    args
		wantErr error
		newComm *Index
	}{
		{
			name:    "change was already added",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Update.copy()}},
			args:    args{chg: gChanges.Foo.None.copy()},
			wantErr: errDuplicatedChg,
		},
		{
			name:    "both identical untracked changes",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Create.copy()}},
			args:    args{chg: gChanges.Foo.Create.copy()},
			newComm: &Index{Changes: []*Change{gChanges.Foo.Create, gChanges.Foo.Create}},
			wantErr: nil,
		},
		{
			name:    "table inconsistency",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Update}},
			args:    args{chg: gChanges.Inconsistent.TableName.copy()},
			wantErr: errNilTable,
		},
		{
			name:    "change modifies the value of a change already in the index",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Update}},
			args:    args{chg: gChanges.Foo.StrValue.copy()},
			newComm: &Index{Changes: []*Change{gChanges.Foo.StrValue.copy().changeType("update")}},
			wantErr: nil,
		},
		{
			name:    "change modifies different col of same schema",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Update}},
			args:    args{chg: gChanges.Foo.ColumnName.copy()},
			newComm: &Index{Changes: []*Change{gChanges.Foo.Update, gChanges.Foo.ColumnName.copy().changeType("update")}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			oldComm := tt.idx
			err := tt.idx.Add(tt.args.chg)
			if err != tt.wantErr {
				t.Errorf("Index.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if diff := cmp.Diff(oldComm, tt.idx); diff != "" {
					t.Errorf("Index.Add() mismatch (-want +got): %s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.newComm, tt.idx); diff != "" {
				t.Errorf("Index.Add() mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestIndex_Rm(t *testing.T) {
	t.Parallel()
	type args struct {
		chg *Change
	}
	tests := []struct {
		name    string
		idx     *Index
		args    args
		wantErr bool
	}{
		{
			name:    "given change doesn't belongs to the index",
			idx:     &Index{Changes: []*Change{gChanges.Foo.None.copy()}},
			args:    args{chg: gChanges.Bar.None.copy()},
			wantErr: true,
		},
		{
			name:    "successfully remove",
			idx:     &Index{Changes: []*Change{gChanges.Foo.None.copy()}},
			args:    args{chg: gChanges.Foo.None.copy()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.idx.Rm(tt.args.chg); (err != nil) != tt.wantErr {
				t.Errorf("Index.Rm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
