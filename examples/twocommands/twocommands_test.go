package main

import (
	"os"
	"testing"

	"github.com/marco-m/rosina"
)

func TestFoo(t *testing.T) {
	want := `hello from FooCmd Run
&main.fooCmd{soft:false}
`
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"foo"})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestBar(t *testing.T) {
	want := `hello from BarCmd Run
&main.barCmd{hard:false}
`
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"bar"})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}
