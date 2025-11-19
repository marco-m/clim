package main

import (
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(MainInt(os.Args[1:]))
}

func MainInt(args []string) int {
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
	cli, err := clim.NewTop("nested", "two subcommands and one nested")
	if err != nil {
		return err
	}

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.Bool(&app.verbose, false), Long: "verbose",
			Help: "Be more verbose",
		}); err != nil {
		return err
	}

	fooCmd, err := newFooCLI(cli)
	if err != nil {
		return err
	}
	barCmd, err := newBarCLI(cli)
	if err != nil {
		return err
	}

	command, err := cli.Parse(args)
	if err != nil {
		return err
	}
	switch command {
	case "nested foo":
		return fooCmd.Run(app)
	case "nested bar list":
		return barCmd.barListCmd.Run(app)
	case "nested bar move":
		return barCmd.barMoveCmd.Run(app)
	default:
		return fmt.Errorf("internal error: unwired command: %s", command)
	}
}
