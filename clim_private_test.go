package clim

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marco-m/rosina"
)

func TestUsage(t *testing.T) {
	type Args struct {
		count  int
		wall   string
		dryRun bool
	}
	var args Args
	cli := New[any]("bang", "bangs head against wall", nil)
	cli.AddFlag(&Flag{
		Value: Int(&args.count, 3),
		Short: "c", Long: "count", Label: "N", Help: "How many times",
	})
	cli.AddFlag(&Flag{
		Value: String(&args.wall, "cardboard"),
		Long:  "wall", Help: "Type of wall",
	})
	cli.AddFlag(&Flag{
		Value: Bool(&args.dryRun, false),
		Long:  "dry-run", Help: "Enable dry-run",
	})

	want := `bang -- bangs head against wall

Usage: bang [options]

Options:

 -c, --count N    How many times (default: 3)
 --dry-run        Enable dry-run (default: false)
 --wall WALL      Type of wall (default: cardboard)

 -h, --help       Print this help and exit
`

	err := cli.usage()

	rosina.AssertErrorIs(t, err, ErrHelp)

	if x := cmp.Diff(want, err.Error()); x != "" {
		t.Fatal("\nwant ---\nhave +++\n", x)
	}
}
