package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/internal/test/thelper"
)

func TestIdx_add(t *testing.T) {
	t.Parallel()
	type args struct {
		chg *Change
	}
	tests := []struct {
		name    string
		idx     *Index
		args    args
		wantErr error
		newIdx  *Index
	}{
		{
			name:    "change was already added",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Update.copy()}},
			args:    args{chg: gChanges.Foo.None.copy()},
			newIdx:  &Index{Changes: []*Change{gChanges.Foo.Update}},
			wantErr: nil,
		},
		{
			name:    "both identical untracked changes",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Create.copy()}},
			args:    args{chg: gChanges.Foo.Create.copy().rmIdAndReturn()},
			newIdx:  &Index{Changes: []*Change{gChanges.Foo.Create, gChanges.Foo.Create.copy().rmIdAndReturn()}},
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
			args:    args{chg: gChanges.Foo.StringValue.copy()},
			newIdx:  &Index{Changes: []*Change{gChanges.Foo.StringValue.copy().changeType("update")}},
			wantErr: nil,
		},
		{
			name:    "change modifies different col of same schema",
			idx:     &Index{Changes: []*Change{gChanges.Foo.Update}},
			args:    args{chg: gChanges.Foo.ColumnName.copy()},
			newIdx:  &Index{Changes: []*Change{gChanges.Foo.Update, gChanges.Foo.ColumnName.copy().changeType("update")}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			oldIdx := tt.idx
			err := tt.idx.add(tt.args.chg)
			if err != tt.wantErr {
				t.Errorf("Index.add() error = %v, wantErr %v", err, tt.wantErr)
			}
			thelper.CmpIfErr(t, err, oldIdx, tt.idx, tt.newIdx, "Index.add()")
		})
	}
}

func TestIndex_rm(t *testing.T) {
	t.Parallel()
	type args struct {
		chg *Change
	}
	tests := []struct {
		name    string
		idx     *Index
		args    args
		wantIdx *Index
	}{
		{
			name:    "given change doesn't belongs to the index",
			idx:     &Index{Changes: []*Change{gChanges.Foo.None.copy()}},
			args:    args{chg: gChanges.Bar.None.copy()},
			wantIdx: &Index{Changes: []*Change{gChanges.Foo.None.copy()}},
		},
		{
			name:    "successfully remove",
			idx:     &Index{Changes: []*Change{gChanges.Foo.None.copy()}},
			args:    args{chg: gChanges.Foo.None.copy()},
			wantIdx: &Index{Changes: []*Change{}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.idx.rm(tt.args.chg)
			if diff := cmp.Diff(tt.wantIdx, tt.idx); diff != "" {
				t.Errorf("Index.rm() mismatch (-want, +got): %s", diff)
			}
		})
	}
}
func (chg *Change) rmIdAndReturn() *Change {
	chg.Id = 0
	return chg
}
