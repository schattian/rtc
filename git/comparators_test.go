package git

import "testing"

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
				chg:      gChanges.Regular.None,
				otherChg: gChanges.Rare.TableName,
			},
			want: false,
		},
		{
			name: "all different",
			args: args{
				chg:      gChanges.Rare.None,
				otherChg: gChanges.Regular.None,
			},
			want: false,
		},
		{
			name: "diff tableName and same entity",
			args: args{
				chg:      gChanges.Regular.None,
				otherChg: gChanges.Regular.TableName,
			},
			want: false,
		},
		{
			name: "both nil entities_id and same tableName",
			args: args{
				chg:      gChanges.Regular.Create,
				otherChg: gChanges.Regular.Create,
			},
			want: false,
		},
		{
			name: "is mirrored",
			args: args{
				chg:      gChanges.Regular.None,
				otherChg: gChanges.Regular.None,
			},
			want: true,
		},
		{
			name: "is mirrored but with diff colName",
			args: args{
				chg:      gChanges.Regular.None,
				otherChg: gChanges.Regular.ColumnName,
			},
			want: true,
		},
		{
			name: "is diff type",
			args: args{
				chg:      gChanges.Regular.Update,
				otherChg: randChg(gChanges.Regular.Retrieve, gChanges.Regular.Delete, gChanges.Regular.Create),
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
			}
		})
	}
}
