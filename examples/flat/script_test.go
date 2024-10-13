// There are two ways to test a program using clim. Pick the one you prefer:
//
// 1. Directly with Go tests (see flat_test.go).
// 2. With rogpeppe/go-internal/testscript (this file and all the
//    .txt files below the testdata/ directory).
//    If you are not familiar with testscript, try it. It is both powerful
//    and simple.
//    For a gentle introduction, see the series starting at
//    https://bitfieldconsulting.com/posts/test-scripts

package main_test

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	flat "github.com/marco-m/clim/examples/flat"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"flat": flat.MainInt,
	}))
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:                 "testdata",
		RequireExplicitExec: true,
	})
}
