package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

/*
$ hg init h
hg init [-e CMD] [--remotecmd CMD] [DEST]

create a new repository in the given directory

    Initialize a new repository in the given directory. If the given directory
    does not exist, it will be created.
...

options:

    --remotecmd CMD specify hg command to run on the remote side
    --mq            operate on patch repository
*/

type initCmd struct {
	remoteCmd string
	mq        bool
}

func newInitCLI(parentCli *clim.CLI) *clim.CLI {
	cli := parentCli.AddCLI("init",
		"create a new repository in the given directory")

	initCmd := initCmd{}

	cli.AddFlag(&clim.Flag{Value: clim.String(&initCmd.remoteCmd, ""),
		Long: "remotecmd", Label: "CMD",
		Desc: "specify hg command to run on the remote side"})
	cli.AddFlag(&clim.Flag{Value: clim.Bool(&initCmd.mq, false),
		Long: "mq", Desc: "operate on patch repository"})

	cli.SetAction(func() error { return initCmd.Run() })

	return cli
}

func (cmd *initCmd) Run() error {
	fmt.Println("hello from InitCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
