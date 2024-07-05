package main

import (
	"os"
	"testing"

	"github.com/go-quicktest/qt"
	"github.com/marco-m/clim/testutils"
)

func TestClone(t *testing.T) {
	want := `hello from CloneCmd Run
&main.cloneCmd{noUpdate:false, updateRev:""}
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"clone"})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}

func TestInit(t *testing.T) {
	want := `hello from InitCmd Run
&main.initCmd{remoteCmd:"", mq:false}
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"init"})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}

func TestIncoming(t *testing.T) {
	want := `hello from IncomingCmd Run
&main.incomingCmd{force:false, newestFirst:false, bundle:"", rev:[]string(nil)}
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"incoming"})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}

func TestOutgoing(t *testing.T) {
	want := `hello from OutgoingCmd Run
&main.outgoingCmd{force:false, rev:[]string(nil), newestFirst:false, bookmarks:false}
`
	reset := testutils.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{"outgoing"})
	qt.Assert(t, qt.IsNil(err))

	out := reset()
	qt.Assert(t, qt.Equals(out, want))
}
