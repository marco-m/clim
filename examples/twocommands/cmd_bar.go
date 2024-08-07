package main

import (
	"context"
	"fmt"

	"github.com/marco-m/clim"
)

type barCmd struct {
	hard bool
}

func newBarCLI(parentCli *clim.CLI[App]) *clim.CLI[App] {
	barCmd := barCmd{}

	cli := parentCli.AddCLI("bar",
		"simple bars all night",
		barCmd.Run)

	cli.AddFlag(&clim.Flag{Value: clim.Bool(&barCmd.hard, false),
		Long: "hard", Desc: "make harder bars"})

	return cli
}

func (cmd *barCmd) Run(ctx context.Context, app App) error {
	fmt.Println("hello from BarCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
