// Program hg shows how to use subcommands with clim, by mimiking a subset of
// the commands of the wonderful mercurial DVCS.

package main

import (
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(MainInt(os.Args[1:]))
}

func MainInt(args []string) int {
	if err := mainErr(args); err != nil {
		fmt.Println(err)
		return clim.ExitCode(err)
	}
	return 0
}

func mainErr(args []string) error {
	cli, err := clim.NewTop("hg", "Mercurial Distributed SCM")
	if err != nil {
		return err
	}

	cloneCmd, err := newCloneCLI(cli)
	if err != nil {
		return err
	}
	initCmd, err := newInitCLI(cli)
	if err != nil {
		return err
	}
	if err := cli.AddGroup("Repository creation",
		cloneCmd.cli, initCmd.cli); err != nil {
		return err
	}

	incomingCmd, err := newIncomingCLI(cli)
	if err != nil {
		return err
	}
	outgoingCmd, err := newOutgoingCLI(cli)
	if err != nil {
		return err
	}
	if err := cli.AddGroup("Remote repository management",
		incomingCmd.cli, outgoingCmd.cli); err != nil {
		return err
	}

	command, err := cli.Parse(args)
	if err != nil {
		return err
	}

	switch command {
	case "hg clone":
		return cloneCmd.Run()
	case "hg init":
		return initCmd.Run()
	case "hg incoming":
		return incomingCmd.Run()
	case "hg outgoing":
		return outgoingCmd.Run()
	default:
		return fmt.Errorf("internal error: unwired command: %s", command)
	}
}
