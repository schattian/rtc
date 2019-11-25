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
				chg:      gChanges.Foo.None,
				otherChg: gChanges.Bar.TableName,
			},
			want: false,
		},
		{
			name: "ALL diff",
			args: args{
				chg:      gChanges.Bar.None,
				otherChg: gChanges.Foo.None,
			},
			want: false,
		},
		{
			name: "DIFF TABLE and SAME ENTITY",
			args: args{
				chg:      gChanges.Foo.None,
				otherChg: gChanges.Foo.TableName,
			},
			want: false,
		},
		{
			name: "NIL ENTITIES and SAME TABLE",
			args: args{
				chg:      gChanges.Foo.Create,
				otherChg: gChanges.Foo.Create,
			},
			want: false,
		},
		{
			name: "but diff TYPE",
			args: args{
				chg:      gChanges.Foo.Update,
				otherChg: randChg(gChanges.Foo.Retrieve, gChanges.Foo.Delete, gChanges.Foo.Create),
			},
			want: false,
		},
		{
			name: "but diff OPTIONS",
			args: args{
				chg:      gChanges.Foo.None,
				otherChg: gChanges.Foo.Options,
			},
			want: false,
		},
		//
		{
			name: "is MIRRORED",
			args: args{
				chg:      gChanges.Foo.None,
				otherChg: gChanges.Foo.None,
			},
			want: true,
		},
		{
			name: "but with diff COLUMN",
			args: args{
				chg:      gChanges.Foo.None,
				otherChg: gChanges.Foo.ColumnName,
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
