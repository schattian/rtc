package fabric

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sebach1/git-crud/integrity"

	"github.com/sebach1/git-crud/internal/test/thelpers"

	"github.com/spf13/afero"
)

func TestFabric_Produce(t *testing.T) {
	t.Parallel()
	type args struct {
		marshal string
	}
	customDir := "testF/testSf"
	tests := []struct {
		name     string
		fabric   *Fabric
		args     args
		wantDir  string
		wantsErr bool
		product  map[integrity.TableName]string // Maps the golden filenames to created goldenFiles
	}{
		{
			name:     "correct usage",
			fabric:   &Fabric{Schema: gSchemas.FooBar},
			args:     args{marshal: "json"},
			wantDir:  fmt.Sprintf("fabric/%v", strings.ToLower(string(gSchemas.FooBar.Name))),
			wantsErr: false,
			product: map[integrity.TableName]string{
				gTables.Foo.Name:    "foo.go",
				gTables.Bar.Name:    "bar.go",
				gTables.FooBar.Name: "foo_bar.go",
			},
		},
		{
			name:     "but w/PRESET DIR",
			fabric:   &Fabric{Schema: gSchemas.FooBar, Dir: customDir},
			args:     args{marshal: "json"},
			wantDir:  customDir,
			wantsErr: false,
			product: map[integrity.TableName]string{
				gTables.Foo.Name:    "foo.go",
				gTables.Bar.Name:    "bar.go",
				gTables.FooBar.Name: "foo_bar.go",
			},
		},
		{
			name:     "SCHEMA does NOT PASS THE VALIdATIONS (is nil)",
			fabric:   &Fabric{Dir: customDir}, // customDir: see that checking os existance of "" dir always returns true
			args:     args{marshal: "json"},
			wantDir:  customDir,
			wantsErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mFs := afero.NewMemMapFs()
			err := tt.fabric.Produce(tt.args.marshal, mFs)
			isCreated := thelpers.IOExist(t, mFs, tt.fabric.Dir, afero.DirExists)
			if tt.fabric.Dir != tt.wantDir {
				t.Errorf("Fabric.Produce() DIFF DIR than expected; want: %v, got: %v", tt.wantDir, tt.fabric.Dir)
			}
			if (err != nil) != tt.wantsErr {
				t.Errorf("Fabric.Produce() err behavior mismatch; wantsErr: %v, got: %v", tt.wantsErr, err)
			}

			if err != nil {
				if isCreated {
					t.Error("Fabric.Produce() GENERATED the DIRectory. EXPECTed to NOT generate it")
				}
				return
			}

			if !isCreated {
				t.Error("Fabric.Produce() did NOT GENERATE the DIRectory. EXPECTed TO generate it")
				return
			}
			for _, table := range tt.fabric.Schema.Blueprint {
				generatedFilename := strings.ToLower(string(table.Name)) + ".go"
				got := thelpers.IOReadFile(t, mFs, tt.fabric.Dir+"/"+generatedFilename)
				thelpers.CmpWithGoldenFile(t, got, fmt.Sprintf("fabric/%s", tt.product[table.Name]), "Fabric.Produce()")
			}
		})
	}
}
