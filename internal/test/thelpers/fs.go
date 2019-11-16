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

// AddFileToFs creates a file with the given filename filled with the content on the baseFs
func AddFileToFs(t *testing.T, filename string, content []byte, baseFs afero.Fs) {
	t.Helper()
	osErr := afero.WriteFile(baseFs, filename, content, os.ModePerm)
	if osErr != nil {
		t.Fatalf("the afero Fs couldnt create the file: %v", osErr)
	}
	return
}

// AddFileToFsByName looks for the filename given over the related testdata dir and creates the file on the baseFs
// Notice the new file will be located in /{filename} path
func AddFileToFsByName(t *testing.T, filename, subset string, baseFs afero.Fs) {
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

	AddFileToFs(t, filename, content, baseFs)
}

// IOExist wraps afero utility func for checking existances (.DirExists, .Exists)
// over the given args and handles the error given of given T
func IOExist(t *testing.T, Fs afero.Fs, sth string, existFunc func(afero.Fs, string) (bool, error)) bool {
	t.Helper()
	res, osErr := existFunc(Fs, sth)
	if osErr != nil {
		t.Fatalf("got UNEXPECTED ERR when trying to use aferos' util: %v", osErr)
	}
	return res
}

// IOReadFile wraps afero.ReadFile
// over the given args and handles the error given of given T
func IOReadFile(t *testing.T, Fs afero.Fs, sth string) []byte {
	t.Helper()
	res, osErr := afero.ReadFile(Fs, sth)
	if osErr != nil {
		t.Fatalf("got UNEXPECTED ERR when trying to use aferos' util: %v", osErr)
	}
	return res
}
