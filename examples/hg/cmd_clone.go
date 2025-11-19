package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type cloneCmd struct {
	noUpdate  bool
	updateRev string
	//
	cli *clim.CLI
}

func newCloneCLI(parent *clim.CLI) (*cloneCmd, error) {
	cloneCmd := &cloneCmd{}

	cli, err := clim.NewSub(parent, "clone",
		"make a copy of an existing repository")
	if err != nil {
		return nil, err
	}
	cloneCmd.cli = cli

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.Bool(&cloneCmd.noUpdate, false),
			Short: "U", Long: "noupdate",
			Help: "the clone will include an empty working directory (only a repository)",
		},
		&clim.Flag{
			Value: clim.String(&cloneCmd.updateRev, ""),
			Short: "u", Long: "updaterev", Label: "REV",
			Help: "revision, tag, or branch to check out",
		}); err != nil {
		return nil, err
	}

	return cloneCmd, nil
}

func (cmd *cloneCmd) Run() error {
	fmt.Println("hello from CloneCmd Run")
	return nil
}
