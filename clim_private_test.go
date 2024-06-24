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
	cli.AddFlag(&Flag{Value: IntVal(&args.count, 3),
		Short: "c", Long: "count", Label: "N", Desc: "How many times"})
	cli.AddFlag(&Flag{Value: StringVal(&args.wall, "cardboard"),
		Long: "wall", Desc: "Type of wall"})
	cli.AddFlag(&Flag{Value: BoolVal(&args.dryRun, false),
		Long: "dry-run", Desc: "Enable dry-run"})

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
