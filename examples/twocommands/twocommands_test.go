package main

import (
	"os"
	"testing"

	"github.com/go-quicktest/qt"
	"github.com/marco-m/clim/testutils"
)

func TestFoo(t *testing.T) {
	want := `hello from FooCmd Run
&main.fooCmd{soft:false}
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"foo"})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}
func TestBar(t *testing.T) {
	want := `hello from BarCmd Run
&main.barCmd{hard:false}
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"bar"})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}
