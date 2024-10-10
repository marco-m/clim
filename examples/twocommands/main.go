package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(mainInt())
}

func mainInt() int {
	err := mainErr(os.Args[1:])
	if err == nil {
		return 0
	}
	fmt.Println(err)
	if errors.Is(err, clim.ErrHelp) {
		return 0
	}
	if errors.Is(err, clim.ErrParse) {
		return 2
	}
	return 1
}

type App struct {
	verbose bool
}

func mainErr(args []string) error {
	app := App{}
	cli, err := clim.New[App]("twocommands", "two simple commands, no groups", nil)
	if err != nil {
		return err
	}

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.Bool(&app.verbose, false),
			Long:  "verbose", Help: "Be more verbose",
		}); err != nil {
		return err
	}

	if err := newFooCLI(cli); err != nil {
		return err
	}
	if err := newBarCLI(cli); err != nil {
		return err
	}

	action, err := cli.Parse(args)
	if err != nil {
		return err
	}

	return action(app)
}
