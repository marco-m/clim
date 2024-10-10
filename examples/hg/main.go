// Program hg shows how to use subcommands with clim, by mimiking a subset of
// the commands of the wonderful mercurial DVCS.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(mainInt())
}

func mainInt() int {
	err := mainErr(os.Args[1:])
	if err == nil {
		return 0
	}
	fmt.Println(err)
	if errors.Is(err, clim.ErrHelp) {
		return 0
	}
	if errors.Is(err, clim.ErrParse) {
		return 2
	}
	return 1
}

type user struct{}

func mainErr(args []string) error {
	cli, err := clim.New[user]("hg", "Mercurial Distributed SCM", nil)
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
