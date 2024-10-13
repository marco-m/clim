package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type cloneCmd struct {
	noUpdate  bool
	updateRev string
}

func newCloneCLI(parent *clim.CLI[user]) (*clim.CLI[user], error) {
	cloneCmd := cloneCmd{}

	cli, err := clim.New(parent, "clone",
		"make a copy of an existing repository",
		cloneCmd.Run)
	if err != nil {
		return nil, err
	}

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

	return cli, nil
}

func (cmd *cloneCmd) Run(uctx user) error {
	fmt.Println("hello from CloneCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
