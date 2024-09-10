// There are two ways to test a program using clim. Pick the one you prefer:
//
// 1. Directly with Go tests (this file).
// 2. With the testscript approach (see script_test.go).

package main

import (
	"os"
	"testing"

	"github.com/marco-m/rosina"
)

func TestBangRun(t *testing.T) {
	want := `1 bang against cardboard
2 bang against cardboard
3 bang against cardboard
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestBangCliHelp(t *testing.T) {
	want := `bang -- bangs head against wall

 Long description.
 Could be multi-line.

Usage: bang [options]

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
	rosina.AssertErrorContains(t, err, want)
}

func TestBangCliWrongInvocation(t *testing.T) {
	want := `unrecognized flag "--foobar"`
	err := mainErr([]string{"--foobar"})
	rosina.AssertErrorContains(t, err, want)
}
