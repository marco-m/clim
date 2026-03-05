// Program hg shows how to use subcommands with clim, by mimiking a subset of
// the commands of the wonderful mercurial DVCS.

package main

import (
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(mainInt(os.Args[1:]))
}

func mainInt(args []string) int {
	if err := mainErr(args); err != nil {
		fmt.Println(err)
		return clim.ExitCode(err)
	}
	return 0
}

type user struct{}

func mainErr(args []string) error {
	cli, err := clim.NewTop[user]("hg", "Mercurial Distributed SCM", nil)
	if err != nil {
		return err
	}

	clonecli, err := newCloneCLI(cli)
	if err != nil {
		return err
	}
	initcli, err := newInitCLI(cli)
	if err != nil {
		return err
	}
	if err := cli.AddGroup("Repository creation",
		clonecli, initcli); err != nil {
		return err
	}

	incomingcli, err := newIncomingCLI(cli)
	if err != nil {
		return err
	}
	outgoingcli, err := newOutgoingCLI(cli)
	if err != nil {
		return err
	}
	if err := cli.AddGroup("Remote repository management",
		incomingcli, outgoingcli); err != nil {
		return err
	}

	action, err := cli.Parse(args)
	if err != nil {
		return err
	}

	uctx := user{}
	return action(uctx)
}
