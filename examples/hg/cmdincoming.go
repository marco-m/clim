package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

/*
hg incoming [-p] [-n] [-M] [-f] [-r REV]... [--bundle FILENAME] [SOURCE]

aliases: in

show new changesets found in source

    Show new changesets found in the specified path/URL or the default pull
    location. These are the changesets that would have been pulled by 'hg
    pull' at the time you issued this command.

    See pull for valid source format details.

    Returns 0 if there are incoming changes, 1 otherwise.

options ([+] can be repeated):

 -f --force             run even if remote repository is unrelated
 -n --newest-first      show newest record first
    --bundle FILE       file to store the bundles into
 -r --rev REV [+]       a remote changeset intended to be added
 ...
*/

type incomingCmd struct {
	force       bool
	newestFirst bool
	bundle      string
	rev         string // FIXME support slices!!!
}

func newIncomingCli(parentCli *clim.Command) *clim.Command {
	cli := parentCli.AddParser("incoming",
		"show new changesets found in source")

	incomingCmd := incomingCmd{}

	cli.AddFlag(clim.BoolVal(&incomingCmd.force, false),
		"f", "force", "",
		"run even if remote repository is unrelated")
	cli.AddFlag(clim.BoolVal(&incomingCmd.newestFirst, false),
		"n", "newest-first", "", "show newest record first")
	cli.AddFlag(clim.StringVal(&incomingCmd.bundle, ""),
		"", "bundle", "FILE",
		"file to store the bundles into")
	cli.AddFlag(clim.StringVal(&incomingCmd.rev, ""),
		"r", "rev", "REV",
		"a remote changeset intended to be added")

	cli.Action(func() error { return incomingCmd.Run() })

	return cli
}

func (cmd *incomingCmd) Run() error {
	fmt.Println("hello from IncomingCmd Run")
	fmt.Printf("%#+v", cmd)
	return nil
}
