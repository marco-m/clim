package main

import (
	"os"
	"testing"

	"github.com/marco-m/rosina"
)

func TestClone(t *testing.T) {
	want := `hello from CloneCmd Run
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"clone"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestInit(t *testing.T) {
	want := `hello from InitCmd Run
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"init"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestIncoming(t *testing.T) {
	want := `hello from IncomingCmd Run
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"incoming"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestOutgoing(t *testing.T) {
	want := `hello from OutgoingCmd Run
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"outgoing"})
	rosina.AssertNoError(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}
