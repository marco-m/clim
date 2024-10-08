package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type fooCmd struct {
	soft        bool
	positionals []string
}

func newFooCLI(parent *clim.CLI[App]) error {
	fooCmd := fooCmd{}

	cli := clim.New("foo", "simple foos all day", fooCmd.Run)

	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&fooCmd.soft, false),
		Long:  "soft", Help: "make softer foos",
	})

	if err := cli.AddPosArgs(&fooCmd.positionals,
		clim.Pair{"COUNT", "How many foos (required)"},
		clim.Pair{"NAME", "Name of the foos (required)"},
		clim.Pair{"COLOR...", "One or more colors (required)"}); err != nil {
		return err
	}

	parent.AddCLI(cli)
	return nil
}

func (cmd *fooCmd) Run(app App) error {
	fmt.Println("hello from FooCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
