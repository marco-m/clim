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
	rootCLI := clim.New("hg", "Mercurial Distributed SCM")

	//
	cloneCLI := newCloneCLI(rootCLI)
	initCLI := newInitCLI(rootCLI)
	rootCLI.AddGroup("Repository creation", cloneCLI, initCLI)

	//
	incomingCLI := newIncomingCLI(rootCLI)
	outgoingCLI := newOutgoingCLI(rootCLI)
	rootCLI.AddGroup("Remote repository management", incomingCLI, outgoingCLI)

	action, err := rootCLI.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return action()
}
