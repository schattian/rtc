package git

import (
	"reflect"
	"testing"

	"github.com/sebach1/git-crud/internal/integrity"
)

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

func TestChange_classifyType(t *testing.T) {
	tests := []struct {
		name    string
		chg     *Change
		want    integrity.CRUD
		wantErr bool
	}{
		{
			name:    "correctly typed CREATE",
			chg:     gChanges.Regular.Create,
			want:    "create",
			wantErr: false,
		},
		{
			name:    "correctly typed RETRIEVE",
			chg:     gChanges.Regular.Retrieve,
			want:    "retrieve",
			wantErr: false,
		},
		{
			name:    "correctly typed UPDATE",
			chg:     gChanges.Regular.Update,
			want:    "update",
			wantErr: false,
		},
		{
			name:    "correctly typed DELETE",
			chg:     gChanges.Regular.Delete,
			want:    "delete",
			wantErr: false,
		},
		{
			name: "unclassifiable inconsistency",
			chg: randChg(gChanges.Inconsistent.Create, gChanges.Inconsistent.Update,
				gChanges.Inconsistent.Delete, gChanges.Inconsistent.Retrieve),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chg.classifyType()
			if (err != nil) != tt.wantErr {
				t.Errorf("Change.classifyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Error(tt.chg.ValueType)
				t.Errorf("Change.classifyType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_validateType(t *testing.T) {
	tests := []struct {
		name    string
		chg     *Change
		wantErr bool
	}{
		{
			name: "unclassifiable inconsistency",
			chg: randChg(gChanges.Inconsistent.Create, gChanges.Inconsistent.Update,
				gChanges.Inconsistent.Delete, gChanges.Inconsistent.Retrieve),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.chg.validateType(); (err != nil) != tt.wantErr {
				t.Errorf("Change.validateType() %v error = %v, wantErr %v", tt.chg.Type, err, tt.wantErr)
			}
		})
	}
}
