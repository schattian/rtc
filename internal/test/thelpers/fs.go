package thelpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func AddFileToFs(t *testing.T, filename string, content []byte, baseFs afero.Fs) afero.Fs {
	t.Helper()
	Fs := afero.NewMemMapFs()
	osErr := afero.WriteFile(Fs, filename, content, os.ModePerm)
	if osErr != nil {
		t.Fatalf("the afero Fs couldnt create the file: %v", osErr)
	}
	return Fs
}

func AddFileToFsByName(t *testing.T, filename string, subset string, baseFs afero.Fs) afero.Fs {
	t.Helper()
	var content []byte
	ext := filepath.Ext(filename)
	testFilename := fmt.Sprintf("testdata/%s", filename)

	var err error
	switch ext {
	case ".jsonnet":
		content, err = exec.Command("jsonnet", testFilename).Output() // Notice avoiding the use of assist pkg due it uses log.Fatal
	default:
		content, err = ioutil.ReadFile(testFilename)
	}
	if err != nil {
		t.Fatalf("the GOLDEN FILE could NOT be READEN: %v", err)
	}

	if subset != "" {
		var set map[string]interface{}
		switch ext {
		case ".jsonnet", ".json":
			err = json.Unmarshal(content, &set)
		}
		if err != nil {
			t.Fatalf("the GOLDEN FILE could NOT be UNMARSHALLED: %v", err)
		}
		content, err = json.Marshal(set[subset])
		if err != nil {
			t.Fatalf("the GOLDEN FILE SUBSET could NOT be MARSHALLED: %v", err)
		}
	}

	return AddFileToFs(t, filename, content, baseFs)
}

func IOutil(t *testing.T, Fs afero.Fs, sth string, utilFunc func(afero.Fs, string) (bool, error)) bool {
	t.Helper()
	res, osErr := utilFunc(Fs, sth)
	if osErr != nil {
		t.Fatalf("got UNEXPECTED ERR when trying to use aferos' util: %v", osErr)
	}
	return res
}
