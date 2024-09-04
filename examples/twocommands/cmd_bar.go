package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type barCmd struct {
	hard bool
}

func newBarCLI() *clim.CLI[App] {
	barCmd := barCmd{}

	cli := clim.New("bar", "simple bars all night", barCmd.Run)

	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&barCmd.hard, false),
		Long:  "hard", Help: "make harder bars",
	})

	return cli
}

func (cmd *barCmd) Run(app App) error {
	fmt.Println("hello from BarCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
