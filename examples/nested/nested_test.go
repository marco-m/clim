package main

import (
	"os"
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestFoo(t *testing.T) {
	want := `hello from FooCmd Run
&main.fooCmd{soft:false, positionals:[]string{}}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"foo"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestBar(t *testing.T) {
	err := mainErr([]string{"bar"})
	rosina.AssertErrorContains(t, err, "expected a command")
}

func TestBarList(t *testing.T) {
	want := `hello from bar list Run
&main.barListCmd{foo:"A"}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"bar", "list", "--foo=A"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertTextEqual(t, out, want, "stdout")
}

func TestBarMove(t *testing.T) {
	want := `hello from bar move Run
&main.barMoveCmd{id:3, dst:"A"}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"bar", "move", "--foo=A", "--id=3"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertTextEqual(t, out, want, "stdout")
}

func TestCliParseError(t *testing.T) {
	err := mainErr([]string{"hello"})
	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `unrecognized command "hello"`)
}
