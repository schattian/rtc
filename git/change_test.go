package git

import (
	"testing"
)

func TestChange_IsUntracked(t *testing.T) {
	tests := []struct {
		name string
		chg  *Change
		want bool
	}{
		{name: "entity_id is not set up", chg: gChanges.Zero, want: true},
		{name: "entity_id is zero-value", chg: gChanges.Regular.Untracked, want: true},

		{name: "entity_id is filled", chg: gChanges.Regular.None, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chg.IsUntracked(); got != tt.want {
				t.Errorf("Change.IsUntracked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCompatibleWith(t *testing.T) {
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
				otherChg: gChanges.Rare.Table,
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
				otherChg: gChanges.Regular.Table,
			},
			want: false,
		},
		{
			name: "both nil entities_id and same tableName",
			args: args{
				chg:      gChanges.Regular.Untracked,
				otherChg: gChanges.Regular.Untracked,
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
				otherChg: gChanges.Regular.Column,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCompatibleWith(tt.args.chg, tt.args.otherChg); got != tt.want {
				t.Errorf("IsCompatibleWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
