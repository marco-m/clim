package main

import (
	"os"
	"testing"

	"github.com/go-quicktest/qt"
	"github.com/marco-m/clim/testutils"
)

func TestBang(t *testing.T) {
	want := `1 bang against cardboard
2 bang against cardboard
3 bang against cardboard
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}
