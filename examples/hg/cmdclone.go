package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type cloneCmd struct {
	noUpdate  bool
	updateRev string
}

func newCloneCLI(parentCli *clim.CLI[user]) *clim.CLI[user] {
	cloneCmd := cloneCmd{}

	cli := parentCli.AddCLI("clone",
		"make a copy of an existing repository",
		cloneCmd.Run)

	cli.AddFlag(&clim.Flag{Value: clim.Bool(&cloneCmd.noUpdate, false),
		Short: "U", Long: "noupdate",
		Desc: "the clone will include an empty working directory (only a repository)"})
	cli.AddFlag(&clim.Flag{Value: clim.String(&cloneCmd.updateRev, ""),
		Short: "u", Long: "updaterev", Label: "REV",
		Desc: "revision, tag, or branch to check out"})

	return cli
}

func (cmd *cloneCmd) Run(uctx user) error {
	fmt.Println("hello from CloneCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
