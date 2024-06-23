package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type cloneCmd struct {
	noUpdate  bool
	updateRev string
}

func newCloneCli(parentCli *clim.Command) *clim.Command {
	cli := parentCli.AddParser("clone",
		"make a copy of an existing repository")

	cloneCmd := cloneCmd{}

	cli.AddFlag(clim.BoolVal(&cloneCmd.noUpdate, false),
		"U", "noupdate", "",
		"the clone will include an empty working directory (only a repository)")
	cli.AddFlag(clim.StringVal(&cloneCmd.updateRev, ""),
		"u", "updaterev", "REV", "revision, tag, or branch to check out")

	cli.Action(func() error { return cloneCmd.Run() })

	return cli
}

func (cmd *cloneCmd) Run() error {
	fmt.Println("hello from CloneCmd Run")
	fmt.Printf("%#+v", cmd)
	return nil
}
