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
	err := mainErr()
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

func mainErr() error {
	rootCli := clim.New("hg", "Mercurial Distributed SCM")

	//
	cloneCli := newCloneCli(rootCli)
	initCmd := newInitCli(rootCli)
	rootCli.Group("Repository creation", cloneCli, initCmd)

	//
	incomingCmd := newIncomingCli(rootCli)
	outgoingCmd := newOutgoingCli(rootCli)
	rootCli.Group("Remote repository management", incomingCmd, outgoingCmd)

	action, err := rootCli.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return action()
}
