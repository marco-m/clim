package main

import (
	"context"
	"fmt"

	"github.com/marco-m/clim"
)

/*
hg outgoing [-M] [-p] [-n] [-f] [-r REV]... [DEST]...

aliases: out

show changesets not found in the destination

    Show changesets not found in the specified destination repository or the
    default push location. These are the changesets that would be pushed if a
    push was requested.

    See pull for details of valid destination formats.

    Returns 0 if there are outgoing changes, 1 otherwise.

options ([+] can be repeated):

 -f --force             run even when the destination is unrelated
 -r --rev REV [+]       a changeset intended to be included in the destination
 -n --newest-first      show newest record first
 -B --bookmarks         compare bookmarks
  ...
*/

type outgoingCmd struct {
	force       bool
	rev         []string
	newestFirst bool
	bookmarks   bool
}

func newOutgoingCLI(parentCli *clim.CLI[user]) *clim.CLI[user] {
	outgoingCmd := outgoingCmd{}

	cli := parentCli.AddCLI("outgoing",
		"show changesets not found in the destination",
		outgoingCmd.Run)

	cli.AddFlag(&clim.Flag{Value: clim.Bool(&outgoingCmd.force, false),
		Short: "f", Long: "force",
		Desc: "run even when the destination is unrelated"})
	cli.AddFlag(&clim.Flag{Value: clim.StringSlice(&outgoingCmd.rev, nil),
		Short: "r", Long: "rev", Label: "REV[,REV,..]",
		Desc: "changeset(s) intended to be included in the destination"})
	cli.AddFlag(&clim.Flag{Value: clim.Bool(&outgoingCmd.newestFirst, false),
		Short: "n", Long: "newest-first", Desc: "show newest record first"})
	cli.AddFlag(&clim.Flag{Value: clim.Bool(&outgoingCmd.bookmarks, false),
		Short: "B", Long: "bookmarks", Desc: "compare bookmarks"})

	return cli
}

func (cmd *outgoingCmd) Run(ctx context.Context, uctx user) error {
	fmt.Println("hello from OutgoingCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
