package schema

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/sebach1/git-crud/internal/test/thelpers"

	"github.com/google/go-cmp/cmp"
)

func TestFromFilename(t *testing.T) {
	t.Parallel()
	decodeValidators := func(sch *Schema) *Schema { sch.applyBuiltinValidators(); return sch } // Skips err checking
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
			want:           decodeValidators(gSchemas.Foo.Copy()),
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
