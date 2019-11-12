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
			name: "DIFF ENTITIES but SAME TABLE",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Rare.TableName,
			},
			want: false,
		},
		{
			name: "ALL diff",
			args: args{
				chg:      gChanges.Rare.None,
				otherChg: gChanges.Basic.None,
			},
			want: false,
		},
		{
			name: "DIFF TABLE and SAME ENTITY",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.TableName,
			},
			want: false,
		},
		{
			name: "NIL ENTITIES and SAME TABLE",
			args: args{
				chg:      gChanges.Basic.Create,
				otherChg: gChanges.Basic.Create,
			},
			want: false,
		},
		{
			name: "but diff TYPE",
			args: args{
				chg:      gChanges.Basic.Update,
				otherChg: randChg(gChanges.Basic.Retrieve, gChanges.Basic.Delete, gChanges.Basic.Create),
			},
			want: false,
		},
		{
			name: "but diff OPTIONS",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.Options,
			},
			want: false,
		},
		//
		{
			name: "is MIRRORED",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.None,
			},
			want: true,
		},
		{
			name: "but with diff COLUMN",
			args: args{
				chg:      gChanges.Basic.None,
				otherChg: gChanges.Basic.ColumnName,
			},
			want: true,
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
