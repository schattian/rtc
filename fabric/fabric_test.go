package fabric

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestFabric_Produce(t *testing.T) {
	t.Parallel()
	type args struct {
		marshal string
	}
	customDir := "testF/testSf"
	tests := []struct {
		name    string
		fabric  *Fabric
		args    args
		wantDir string
		wantErr bool
	}{
		{
			name:    "correct usage",
			fabric:  &Fabric{Schema: gSchemas.Basic},
			args:    args{marshal: "json"},
			wantDir: fmt.Sprintf("fabric/%v", strings.ToLower(string(gSchemas.Basic.Name))),
			wantErr: false,
		},
		{
			name:    "but w/PRESET DIR",
			fabric:  &Fabric{Schema: gSchemas.Basic, Dir: customDir},
			args:    args{marshal: "json"},
			wantDir: customDir,
			wantErr: false,
		},
		{
			name:    "SCHEMA does NOT PASS THE VALIDATIONS (is nil)",
			fabric:  &Fabric{Dir: customDir}, // customDir: see that checking os existance of "" dir always returns true
			args:    args{marshal: "json"},
			wantDir: customDir,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mFs := afero.NewMemMapFs()
			err := tt.fabric.Produce(tt.args.marshal, mFs)
			isCreated, osErr := afero.DirExists(mFs, tt.fabric.Dir)
			if osErr != nil {
				t.Errorf("Fabric.Produce() got UNEXPECTED ERR when trying to fetch dir on aferos' MemMapFs: %v", err)
			}
			if tt.fabric.Dir != tt.wantDir {
				t.Errorf("Fabric.Produce() DIFF DIR than expected; want: %v, got: %v", tt.wantDir, tt.fabric.Dir)
			}

			if err != nil {
				if isCreated {
					t.Error("Fabric.Produce() GENERATED the DIRectory. EXPECTed to NOT generate it")
				}
				return
			}

			if !isCreated {
				t.Error("Fabric.Produce() did NOT GENERATE the DIRectory. EXPECTed TO generate it")
			}
		})
	}
}
