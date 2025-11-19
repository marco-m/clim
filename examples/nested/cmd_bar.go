package main

import (
	"fmt"

	"github.com/marco-m/clim"
)

type barCmds struct {
	*barListCmd
	*barMoveCmd
}

type barListCmd struct {
	foo string
}

type barMoveCmd struct {
	id  int
	dst string
}

func newBarCLI(parent *clim.CLI) (*barCmds, error) {
	barCmds := barCmds{}
	cli, err := clim.NewSub(parent, "bar", "simple bars all night; has subcommands")
	if err != nil {
		return nil, err
	}

	barCmds.barListCmd, err = newBarListCLI(cli)
	if err != nil {
		return nil, err
	}
	barCmds.barMoveCmd, err = newBarMoveCLI(cli)
	if err != nil {
		return nil, err
	}

	return &barCmds, nil
}

//
//
//

func newBarListCLI(parent *clim.CLI) (*barListCmd, error) {
	barListCmd := barListCmd{}

	cli, err := clim.NewSub(parent, "list", "list all bars in a given foo")
	if err != nil {
		return nil, err
	}

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.String(&barListCmd.foo, ""),
			Long:  "foo", Help: "Name of the foo (see nested foo list)",
			Required: true,
		}); err != nil {
		return nil, err
	}

	return &barListCmd, nil
}

func (cmd *barListCmd) Run(app App) error {
	fmt.Println("hello from bar list Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}

//
//
//

func newBarMoveCLI(parent *clim.CLI) (*barMoveCmd, error) {
	barMoveCmd := barMoveCmd{}

	cli, err := clim.NewSub(parent, "move", "move a bar into a foo")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &barMoveCmd, nil
}

func (cmd *barMoveCmd) Run(app App) error {
	fmt.Println("hello from bar move Run")
	fmt.Printf("%#+v\n", cmd)
	return nil
}
