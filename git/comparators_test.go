package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAreCompatible(t *testing.T) {
	t.Parallel()
	type args struct {
		chg      *Change
		otherChg *Change
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "different entities but same tableName",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Rare.TableName,
			},
			want: false,
		},
		{
			name: "all different",
			args: args{
				chg:      gChanges.Rare.None,
				otherChg: gChanges.Basic.None,
			},
			want: false,
		},
		{
			name: "diff tableName and same entity",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.TableName,
			},
			want: false,
		},
		{
			name: "both nil entities_id and same tableName",
			args: args{
				chg:      gChanges.Basic.Create,
				otherChg: gChanges.Basic.Create,
			},
			want: false,
		},
		{
			name: "is mirrored",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.None,
			},
			want: true,
		},
		{
			name: "is mirrored but with diff colName",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.ColumnName,
			},
			want: true,
		},
		{
			name: "is diff type",
			args: args{
				chg:      gChanges.Basic.Update,
				otherChg: randChg(gChanges.Basic.Retrieve, gChanges.Basic.Delete, gChanges.Basic.Create),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := AreCompatible(tt.args.chg, tt.args.otherChg); got != tt.want {
				t.Errorf("AreCompatibleWith() = %v, want %v", got, tt.want)
				diff := cmp.Diff(tt.args.chg, tt.args.otherChg)
				t.Error(diff)
			}
		})
	}
}
