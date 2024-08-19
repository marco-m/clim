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
	cli := clim.New[App]("twocommands", "two simple commands, no groups", nil)

	cli.AddFlag(&clim.Flag{Value: clim.Bool(&app.verbose, false),
		Long: "verbose", Desc: "Be more verbose",
	})

	newFooCLI(cli)
	newBarCLI(cli)

	action, err := cli.Parse(args)
	if err != nil {
		return err
	}

	return action(app)
}
