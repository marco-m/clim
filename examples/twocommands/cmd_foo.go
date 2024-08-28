package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type fooCmd struct {
	soft bool
}

func newFooCLI(parentCli *clim.CLI[App]) *clim.CLI[App] {
	fooCmd := fooCmd{}

	cli := parentCli.AddCLI("foo",
		"simple foos all day",
		fooCmd.Run)

	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&fooCmd.soft, false),
		Long:  "soft", Help: "make softer foos",
	})

	return cli
}

func (cmd *fooCmd) Run(app App) error {
	fmt.Println("hello from FooCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
