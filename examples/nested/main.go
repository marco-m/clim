package main

import (
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(mainInt(os.Args[1:]))
}

func mainInt(args []string) int {
	if err := mainErr(args); err != nil {
		fmt.Println(err)
		return clim.ExitCode(err)
	}
	return 0
}

type App struct {
	verbose bool
}

func mainErr(args []string) error {
	app := App{}
	cli, err := clim.NewTop[App]("nested", "two subcommands and one nested", nil)
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
