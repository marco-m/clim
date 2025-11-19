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
	//
	cli *clim.CLI
}

func newInitCLI(parent *clim.CLI) (*initCmd, error) {
	initCmd := &initCmd{}

	cli, err := clim.NewSub(parent, "init",
		"create a new repository in the given directory")
	if err != nil {
		return nil, err
	}
	initCmd.cli = cli

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.String(&initCmd.remoteCmd, ""),
			Long:  "remotecmd", Label: "CMD",
			Help: "specify hg command to run on the remote side",
		},
		&clim.Flag{
			Value: clim.Bool(&initCmd.mq, false),
			Long:  "mq", Help: "operate on patch repository",
		}); err != nil {
		return nil, err
	}

	return initCmd, nil
}

func (cmd *initCmd) Run() error {
	fmt.Println("hello from InitCmd Run")
	return nil
}
