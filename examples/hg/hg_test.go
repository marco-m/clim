package main

import (
	"os"
	"testing"

	"github.com/marco-m/rosina"
)

func TestClone(t *testing.T) {
	want := `hello from CloneCmd Run
&main.cloneCmd{noUpdate:false, updateRev:""}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"clone"})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestInit(t *testing.T) {
	want := `hello from InitCmd Run
&main.initCmd{remoteCmd:"", mq:false}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"init"})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestIncoming(t *testing.T) {
	want := `hello from IncomingCmd Run
&main.incomingCmd{force:false, newestFirst:false, bundle:"", rev:[]string(nil)}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"incoming"})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestOutgoing(t *testing.T) {
	want := `hello from OutgoingCmd Run
&main.outgoingCmd{force:false, rev:[]string(nil), newestFirst:false, bookmarks:false}
`
	readReset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"outgoing"})
	rosina.AssertIsNil(t, err)

	out := readReset()
	rosina.AssertEqual(t, out, want, "stdout")
}
