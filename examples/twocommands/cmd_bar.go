package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type barCmd struct {
	hard bool
}

func newBarCLI(parent *clim.CLI[App]) error {
	barCmd := barCmd{}

	cli, err := clim.New("bar", "simple bars all night", barCmd.Run)
	if err != nil {
		return err
	}

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.Bool(&barCmd.hard, false),
			Long:  "hard", Help: "make harder bars",
		},
	); err != nil {
		return err
	}

	parent.AddCLI(cli)
	return nil
}

func (cmd *barCmd) Run(app App) error {
	fmt.Println("hello from BarCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
