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

	cli, err := clim.New(parent, "foo", "simple foos all day", fooCmd.Run)
	if err != nil {
		return err
	}

	if err := cli.AddFlags(&clim.Flag{
		Value: clim.Bool(&fooCmd.soft, false),
		Long:  "soft", Help: "make softer foos",
	}); err != nil {
		return err
	}

	if err := cli.AddPosArgs(&fooCmd.positionals,
		clim.Pair{"COUNT", "How many foos (required)"},
		clim.Pair{"NAME", "Name of the foos (required)"},
		clim.Pair{"COLOR...", "One or more colors (required)"}); err != nil {
		return err
	}

	return nil
}

func (cmd *fooCmd) Run(app App) error {
	fmt.Println("hello from FooCmd Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
