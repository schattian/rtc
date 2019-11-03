package thelpers

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CmpIfErr(t *testing.T, err error, old, actual, new interface{}, msg string) {
	if err != nil {
		if diff := cmp.Diff(old, actual); diff != "" {
			t.Errorf("%s errored mismatch (-want +got): %s", msg, diff)
		}
		return
	}
	if diff := cmp.Diff(new, actual); diff != "" {
		t.Errorf("%s mismatch (-want +got): %s", msg, diff)
	}
}

func CmpWithGoldenFile(t *testing.T, got []byte, goldenFilename, msg string) {
	t.Helper()

	want, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", goldenFilename))
	if err != nil {
		t.Fatalf("Error trying to read the goldenFile: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("%s mismatch (-want +got): %s", msg, diff)
	}
}
