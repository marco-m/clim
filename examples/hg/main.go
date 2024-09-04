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
	rootCLI := clim.New[user]("hg", "Mercurial Distributed SCM", nil)

	//
	rootCLI.AddGroup("Repository creation",
		rootCLI.AddCLI(newCloneCLI()),
		rootCLI.AddCLI(newInitCLI()))

	//
	rootCLI.AddGroup("Remote repository management",
		rootCLI.AddCLI(newIncomingCLI()),
		rootCLI.AddCLI(newOutgoingCLI()))

	action, err := rootCLI.Parse(args)
	if err != nil {
		return err
	}

	uctx := user{}
	return action(uctx)
}
