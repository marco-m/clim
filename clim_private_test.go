package clim

import (
	"strings"
	"testing"

	"github.com/go-quicktest/qt"
	"github.com/google/go-cmp/cmp"
)

func TestUsage(t *testing.T) {
	type Args struct {
		count  int
		wall   string
		dryRun bool
	}
	var args Args
	cli := New("bang", "bangs head against wall")
	cli.IntVar(&args.count, "c", "count", "N", 3, "How many times")
	cli.StringVar(&args.wall, "", "wall", "WALL", "cardboard", "Type of wall")
	cli.BoolVar(&args.dryRun, "", "dry-run", false, "Enable dry-run")

	want := strings.TrimSpace(`
bang -- bangs head against wall

Usage: bang [options]

Options:

 -c, --count N    How many times (default: 3)
 --dry-run        Enable dry-run (default: false)
 --wall WALL      Type of wall (default: cardboard)

 -h, --help       Print this help and exit
`)

	err := cli.usage()

	qt.Assert(t, qt.ErrorIs(err, ErrHelp))
	if x := cmp.Diff(want, err.Error()); x != "" {
		t.Fatal("\nwant ---\nhave +++\n", x)
	}
}
