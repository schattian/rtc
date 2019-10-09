package git

import (
	"reflect"
	"testing"
)

func TestAreCompatible(t *testing.T) {
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
				otherChg: gChanges.Regular.Column,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AreCompatible(tt.args.chg, tt.args.otherChg); got != tt.want {
				t.Errorf("AreCompatibleWith() = %v, want %v", got, tt.want)
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
		name          string
		chg           *Change
		args          args
		wantErr       bool
		wantValueType string
	}{
		{
			name:          "string",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Regular.StrValue.StrValue},
			wantErr:       false,
			wantValueType: "string",
		},
		{
			name:          "json",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Regular.JSONValue.JSONValue},
			wantErr:       false,
			wantValueType: "json",
		},
		{
			name:          "int",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Regular.IntValue.IntValue},
			wantErr:       false,
			wantValueType: "int",
		},
		{
			name:          "float32",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Regular.Float32Value.Float32Value},
			wantErr:       false,
			wantValueType: "float32",
		},
		{
			name:          "float64",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Regular.Float64Value.Float64Value},
			wantErr:       false,
			wantValueType: "float64",
		},
		{
			name:          "nil",
			chg:           randChg(cleansedChgs...),
			args:          args{val: nil},
			wantErr:       true,
			wantValueType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.chg.SetValue(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Change.SetValue() error got: %v, wantErr %v", err, tt.wantErr)
			}
			if gotValueType := tt.chg.ValueType; gotValueType != tt.wantValueType {
				t.Errorf("Change.SetValue() type saved mismatch; got: %v, want: %v", gotValueType, tt.wantValueType)
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
