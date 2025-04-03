package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// compare structs with ignored fields
func Diff(t *testing.T, expected any, actual any, ignoredFields ...string) {
	t.Helper()
	ignored := cmpopts.IgnoreFields(expected, ignoredFields...)
	if d := cmp.Diff(expected, actual, ignored); d != "" {
		t.Errorf("mismatch (-want +got):\n%s", d)
	}
}
