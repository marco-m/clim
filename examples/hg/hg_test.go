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
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"clone"})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestInit(t *testing.T) {
	want := `hello from InitCmd Run
&main.initCmd{remoteCmd:"", mq:false}
`
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"init"})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestIncoming(t *testing.T) {
	want := `hello from IncomingCmd Run
&main.incomingCmd{force:false, newestFirst:false, bundle:"", rev:[]string(nil)}
`
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"incoming"})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}

func TestOutgoing(t *testing.T) {
	want := `hello from OutgoingCmd Run
&main.outgoingCmd{force:false, rev:[]string(nil), newestFirst:false, bookmarks:false}
`
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"outgoing"})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}
