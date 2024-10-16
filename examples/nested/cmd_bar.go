package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

func newBarCLI(parent *clim.CLI[App]) error {
	cli, err := clim.NewSub[App](parent, "bar",
		"simple bars all night; has subcommands", nil)
	if err != nil {
		return err
	}

	if err := newBarListCLI(cli); err != nil {
		return err
	}
	if err := newBarMoveCLI(cli); err != nil {
		return err
	}

	return nil
}

//
//
//

type barListCmd struct {
	foo string
}

func newBarListCLI(parent *clim.CLI[App]) error {
	barListCmd := barListCmd{}

	cli, err := clim.NewSub(parent, "list", "list all bars in a given foo",
		barListCmd.Run)
	if err != nil {
		return err
	}

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.String(&barListCmd.foo, ""),
			Long:  "foo", Help: "Name of the foo (see nested foo list)",
			Required: true,
		}); err != nil {
		return err
	}

	return nil
}

func (cmd *barListCmd) Run(app App) error {
	fmt.Println("hello from bar list Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}

//
//
//

type barMoveCmd struct {
	id  int
	dst string
}

func newBarMoveCLI(parent *clim.CLI[App]) error {
	barMoveCmd := barMoveCmd{}

	cli, err := clim.NewSub(parent, "move",
		"move a bar into a foo", barMoveCmd.Run)
	if err != nil {
		return err
	}

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.Int(&barMoveCmd.id, 0),
			Long:  "id", Help: "bar ID",
			Required: true,
		},
		&clim.Flag{
			Value: clim.String(&barMoveCmd.dst, ""),
			Long:  "foo", Help: "Foo name",
			Required: true,
		}); err != nil {
		return err
	}

	return nil
}

func (cmd *barMoveCmd) Run(app App) error {
	fmt.Println("hello from bar move Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
