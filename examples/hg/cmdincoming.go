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
	rev         []string
}

func newIncomingCLI(parentCli *clim.CLI[user]) *clim.CLI[user] {
	incomingCmd := incomingCmd{}

	cli := parentCli.AddCLI("incoming",
		"show new changesets found in source",
		incomingCmd.Run)

	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&incomingCmd.force, false),
		Short: "f", Long: "force",
		Help: "run even if remote repository is unrelated",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&incomingCmd.newestFirst, false),
		Short: "n", Long: "newest-first",
		Help: "show newest record first",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.String(&incomingCmd.bundle, ""),
		Long:  "bundle", Label: "FILE",
		Help: "file to store the bundles into",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.StringSlice(&incomingCmd.rev, nil),
		Short: "r", Long: "rev", Label: "REV[,REV,..]",
		Help: "remote changeset(s) intended to be added",
	})

	return cli
}

func (cmd *incomingCmd) Run(uctx user) error {
	fmt.Println("hello from IncomingCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
