// There are two ways to test a program using clim. Pick the one you prefer:
//
// 1. Directly with Go tests (this file).
// 2. With the testscript approach (see script_test.go).

package main

import (
	"os"
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestFlatRun(t *testing.T) {
	want := `1 flatten against cardboard
2 flatten against cardboard
3 flatten against cardboard
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertTextEqual(t, out, want, "stdout")
}

func TestFlatCliHelp(t *testing.T) {
	want := `flat -- flattens head against wall

 Long description.
 Could be multi-line.

Usage: flat [options]

Examples:

 One or more examples.

 Could be multi-line.

Options:

 -c, --count N         How many times (default: 3)
 --doors N[,N,..]      Doors sequence
 --dry-run             Enable dry-run (default: false)
 --floors F[,F,..]     Floors sequence
 --wall WALL           Type of wall (default: cardboard)
 --windows N[,N,..]    Windows sequence

 -h, --help            Print this help and exit

 For more information, see https://www.example.org/
`
	err := mainErr([]string{"-h"})
	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertTextEqual(t, err.Error(), want, "error message")
}

func TestFlatCliWrongInvocation(t *testing.T) {
	want := `unrecognized flag "--foobar"`
	err := mainErr([]string{"--foobar"})
	rosina.AssertErrorContains(t, err, want)
}
