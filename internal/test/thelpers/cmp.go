package thelpers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CmpIfErr(t *testing.T, err error, old, actual, new interface{}, msg string) {
	if err != nil {
		if diff := cmp.Diff(old, actual); diff != "" {
			t.Errorf("%v errored mismatch (-want +got): %s", msg, diff)
		}
		return
	}
	if diff := cmp.Diff(new, actual); diff != "" {
		t.Errorf("%v mismatch (-want +got): %s", msg, diff)
	}
}
