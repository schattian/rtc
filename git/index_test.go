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
		want    error
		newComm *Index
	}{
		{
			name: "change was already added",
			idx:  &Index{Changes: []*Change{gChanges.Basic.Update.copy()}},
			args: args{chg: gChanges.Basic.None.copy()},
			want: errDuplicatedChg,
		},
		{
			name:    "both identical untracked changes",
			idx:     &Index{Changes: []*Change{gChanges.Basic.Create.copy()}},
			args:    args{chg: gChanges.Basic.Create.copy()},
			newComm: &Index{Changes: []*Change{gChanges.Basic.Create, gChanges.Basic.Create}},
		},
		{
			name: "table inconsistency",
			idx:  &Index{Changes: []*Change{gChanges.Basic.Update}},
			args: args{chg: gChanges.Inconsistent.TableName.copy()},
			want: errNilTable,
		},
		{
			name:    "change modifies the value of a change already in the index",
			idx:     &Index{Changes: []*Change{gChanges.Basic.Update}},
			args:    args{chg: gChanges.Basic.StrValue.copy()},
			newComm: &Index{Changes: []*Change{gChanges.Basic.StrValue.copy().changeType("update")}},
		},
		{
			name:    "change modifies different col of same schema",
			idx:     &Index{Changes: []*Change{gChanges.Basic.Update}},
			args:    args{chg: gChanges.Basic.ColumnName.copy()},
			newComm: &Index{Changes: []*Change{gChanges.Basic.Update, gChanges.Basic.ColumnName.copy().changeType("update")}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			oldComm := tt.idx
			err := tt.idx.Add(tt.args.chg)
			if err != tt.want {
				t.Errorf("Index.Add() error = %v, wantErr %v", err, tt.want)
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
			idx:     &Index{Changes: []*Change{gChanges.Basic.None.copy()}},
			args:    args{chg: gChanges.Rare.None.copy()},
			wantErr: true,
		},
		{
			name:    "successfully remove",
			idx:     &Index{Changes: []*Change{gChanges.Basic.None.copy()}},
			args:    args{chg: gChanges.Basic.None.copy()},
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
