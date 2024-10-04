package main

import (
	"os"
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestFoo(t *testing.T) {
	want := `hello from FooCmd Run
&main.fooCmd{soft:false}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"foo"})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestBar(t *testing.T) {
	want := `hello from BarCmd Run
&main.barCmd{hard:false}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"bar"})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestCliParseError(t *testing.T) {
	err := mainErr([]string{"hello"})
	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `unrecognized command "hello"`)
}
