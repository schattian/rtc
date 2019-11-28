package schema

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/sebach1/rtc/internal/test/thelpers"

	"github.com/google/go-cmp/cmp"
)

func UnmarshalValidatorsAndReturn(t *testing.T, sch *Schema) *Schema {
	t.Helper()
	err := sch.applyBuiltinValidators()
	if err != nil {
		t.Fatalf("Couldn't unmarshal validators at helper layer: %v", err)
	}
	return sch
}

func TestFromFilename(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		goldenFilename string
		fake           bool // Fake the goldenFile (w/ empty content)
		want           *Schema
		wantErr        error
	}{
		{
			name:           "CORRECT USAGE",
			goldenFilename: "schemas.jsonnet",
			want:           UnmarshalValidatorsAndReturn(t, gSchemas.Foo.Copy()),
			wantErr:        nil,
		},
		{
			name:           "the schema contains a COLUMN WITH INCONSISTENT VALUE TYPE", // Ensure err checking in applyBuiltinValidators
			goldenFilename: "inconsistent_schemas.jsonnet",
			wantErr:        errUnallowedColumnType,
		},
		{
			name:           "the EXT is NOT ALLOWED",
			goldenFilename: "schemas.matlab", // (?)
			fake:           true,
			wantErr:        errUnallowedExt,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Fs := afero.NewMemMapFs()
			if tt.fake {
				thelpers.AddFileToFs(t, tt.goldenFilename, []byte{}, Fs)
			} else {
				thelpers.AddFileToFsByName(t, tt.goldenFilename, "foo", Fs)
			}
			got, err := FromFilename(tt.goldenFilename, Fs)
			if err != tt.wantErr {
				t.Errorf("FromFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FromFilename() mismatch (-want +got): %s", diff)
			}
		})
	}
}
