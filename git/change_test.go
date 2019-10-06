package git

import (
	"reflect"
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

func TestChange_SetValue(t *testing.T) {
	type args struct {
		val interface{}
	}
	cleansedChgs := []*Change{gChanges.Regular.CleanValue, gChanges.Rare.CleanValue}
	tests := []struct {
		name     string
		chg      *Change
		args     args
		wantErr  bool
		wantType string
	}{
		{
			name:     "string",
			chg:      randChg(cleansedChgs...),
			args:     args{val: gChanges.Regular.StrValue.StrValue},
			wantErr:  false,
			wantType: "string",
		},
		{
			name:     "json",
			chg:      randChg(cleansedChgs...),
			args:     args{val: gChanges.Regular.JSONValue.JSONValue},
			wantErr:  false,
			wantType: "json",
		},
		{
			name:     "int",
			chg:      randChg(cleansedChgs...),
			args:     args{val: gChanges.Regular.IntValue.IntValue},
			wantErr:  false,
			wantType: "int",
		},
		{
			name:     "float32",
			chg:      randChg(cleansedChgs...),
			args:     args{val: gChanges.Regular.Float32Value.Float32Value},
			wantErr:  false,
			wantType: "float32",
		},
		{
			name:     "float64",
			chg:      randChg(cleansedChgs...),
			args:     args{val: gChanges.Regular.Float64Value.Float64Value},
			wantErr:  false,
			wantType: "float64",
		},
		{
			name:     "nil",
			chg:      randChg(cleansedChgs...),
			args:     args{val: nil},
			wantErr:  true,
			wantType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.chg.SetValue(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Change.SetValue() error got: %v, wantErr %v", err, tt.wantErr)
			}
			if gotType := tt.chg.Type; gotType != tt.wantType {
				t.Errorf("Change.SetValue() type saved mismatch; got: %v, want: %v", gotType, tt.wantType)
			}
		})
	}
}

func TestChange_Value(t *testing.T) {
	tests := []struct {
		name string
		chg  *Change
		want interface{}
	}{
		{
			name: "json",
			chg:  gChanges.Regular.JSONValue,
			want: gChanges.Regular.JSONValue.JSONValue,
		},
		{
			name: "str",
			chg:  gChanges.Regular.None,
			want: gChanges.Regular.None.StrValue,
		},
		{
			name: "float32",
			chg:  gChanges.Regular.Float32Value,
			want: gChanges.Regular.Float32Value.Float32Value,
		},
		{
			name: "int",
			chg:  gChanges.Regular.IntValue,
			want: gChanges.Regular.IntValue.IntValue,
		},
		{
			name: "float64",
			chg:  gChanges.Regular.Float64Value,
			want: gChanges.Regular.Float64Value.Float64Value,
		},
		{
			name: "str",
			chg:  gChanges.Regular.None,
			want: gChanges.Regular.None.StrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chg.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
