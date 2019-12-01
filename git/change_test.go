package git

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/rtc/integrity"
)

func TestChange_SetValue(t *testing.T) {
	t.Parallel()
	type args struct {
		val interface{}
	}
	cleansedChgs := []*Change{gChanges.Foo.CleanValue, gChanges.Bar.CleanValue}
	tests := []struct {
		name          string
		chg           *Change
		args          args
		wantErr       error
		wantValueType integrity.ValueType
	}{
		{
			name:          "string",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Foo.StrValue.StrValue},
			wantErr:       nil,
			wantValueType: "string",
		},
		{
			name:          "json",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Foo.JSONValue.JSONValue},
			wantErr:       nil,
			wantValueType: "json",
		},
		{
			name:          "int",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Foo.IntValue.IntValue},
			wantErr:       nil,
			wantValueType: "int",
		},
		{
			name:          "float32",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Foo.Float32Value.Float32Value},
			wantErr:       nil,
			wantValueType: "float32",
		},
		{
			name:          "float64",
			chg:           randChg(cleansedChgs...),
			args:          args{val: gChanges.Foo.Float64Value.Float64Value},
			wantErr:       nil,
			wantValueType: "float64",
		},
		{
			name:    "nil",
			chg:     randChg(cleansedChgs...),
			args:    args{val: nil},
			wantErr: errUnsafeValueType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.chg.SetValue(tt.args.val); err != tt.wantErr {
				t.Errorf("Change.SetValue() error got: %v, wantErr %v", err, tt.wantErr)
			}
			if gotValueType := tt.chg.ValueType; gotValueType != tt.wantValueType {
				t.Errorf("Change.SetValue() type saved mismatch; got: %v, want: %v", gotValueType, tt.wantValueType)
			}
		})
	}
}

func TestChange_Value(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		chg  *Change
		want interface{}
	}{
		{
			name: "json",
			chg:  gChanges.Foo.JSONValue,
			want: gChanges.Foo.JSONValue.JSONValue,
		},
		{
			name: "str",
			chg:  gChanges.Foo.None,
			want: gChanges.Foo.None.StrValue,
		},
		{
			name: "float32",
			chg:  gChanges.Foo.Float32Value,
			want: gChanges.Foo.Float32Value.Float32Value,
		},
		{
			name: "int",
			chg:  gChanges.Foo.IntValue,
			want: gChanges.Foo.IntValue.IntValue,
		},
		{
			name: "float64",
			chg:  gChanges.Foo.Float64Value,
			want: gChanges.Foo.Float64Value.Float64Value,
		},
		{
			name: "str",
			chg:  gChanges.Foo.None,
			want: gChanges.Foo.None.StrValue,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.chg.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Change.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChange_classifyType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		chg     *Change
		want    integrity.CRUD
		wantErr error
	}{
		{
			name:    "correctly typed CREATE",
			chg:     gChanges.Foo.Create,
			want:    "create",
			wantErr: nil,
		},
		{
			name:    "correctly typed RETRIEVE",
			chg:     gChanges.Foo.Retrieve,
			want:    "retrieve",
			wantErr: nil,
		},
		{
			name:    "correctly typed UPDATE",
			chg:     gChanges.Foo.Update,
			want:    "update",
			wantErr: nil,
		},
		{
			name:    "correctly typed DELETE",
			chg:     gChanges.Foo.Delete,
			want:    "delete",
			wantErr: nil,
		},
		{
			name: "unclassifiable inconsistency",
			chg: randChg(gChanges.Inconsistent.Create, gChanges.Inconsistent.Update,
				gChanges.Inconsistent.Delete, gChanges.Inconsistent.Retrieve),
			wantErr: errUnclassifiableChg,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.chg.classifyType()
			if err != tt.wantErr {
				t.Errorf("Change.classifyType() error = %v, wantErr %v", err, tt.want)
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
	t.Parallel()
	tests := []struct {
		name     string
		chg      *Change
		wantsErr bool
	}{
		{
			name: "unclassifiable inconsistency",
			chg: randChg(gChanges.Inconsistent.Create, gChanges.Inconsistent.Update,
				gChanges.Inconsistent.Delete, gChanges.Inconsistent.Retrieve),
			wantsErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.chg.validateType(); (err != nil) != tt.wantsErr {
				t.Errorf("Change.validateType() %v error = %v, wantsErr %v", tt.chg.Type, err, tt.wantsErr)
			}
		})
	}
}

func TestChange_SetOption(t *testing.T) {
	t.Parallel()
	type args struct {
		key integrity.OptionKey
		val interface{}
	}

	assignAndReturn := func(opts Options, k integrity.OptionKey, v interface{}) Options {
		opts[k] = v
		return opts
	}

	tests := []struct {
		name        string
		chg         *Change
		args        args
		wantErr     error
		wantOptions Options
	}{
		{
			name:        "ALREADY INITIALISED options",
			chg:         gChanges.Foo.None.copy(),
			args:        args{key: "testKey", val: "testVal"},
			wantErr:     nil,
			wantOptions: assignAndReturn(gChanges.Foo.None.copy().Options, "testKey", "testVal"),
		},
		{
			name:        "UNINITIALIZED options",
			chg:         gChanges.Zero.copy(),
			args:        args{key: "testKey", val: "testVal"},
			wantErr:     nil,
			wantOptions: Options{"testKey": "testVal"},
		},
		{
			name:    "EMPTY KEY ERROR",
			chg:     gChanges.Zero.copy(),
			args:    args{key: "", val: "testVal"},
			wantErr: errNilOptionKey,
		},
		{
			name:        "CHANGE OPTION VALUE",
			chg:         gChanges.Foo.None.copy(),
			args:        args{key: gChanges.Foo.None.Options.Keys()[0], val: "testVal"},
			wantErr:     nil,
			wantOptions: Options{gChanges.Foo.None.Options.Keys()[0]: "testVal"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			oldOptions := tt.chg.Options
			err := tt.chg.SetOption(tt.args.key, tt.args.val)
			if err != tt.wantErr {
				t.Errorf("Change.SetOption() error = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				if diff := cmp.Diff(oldOptions, tt.chg.Options); diff != "" {
					t.Errorf("Change.SetOption() mismatch (-want +got): %s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.wantOptions, tt.chg.Options); diff != "" {
				t.Errorf("Change.SetOption() mismatch (-want +got): %s", diff)
			}
		})
	}
}
