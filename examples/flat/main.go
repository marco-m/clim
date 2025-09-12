package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/marco-m/clim"
)

func main() {
	// Alternative:
	// os.Exit(clim.ExitCode(mainErr, os.Args[1:], os.Stderr))
	os.Exit(MainInt())
}

func MainInt() int {
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

type Application struct {
	count   int
	wall    string
	dryRun  bool
	windows []int
	doors   []int
	floors  []string
}

func mainErr(args []string) error {
	var app Application
	cli, err := clim.NewTop("flat", "flattens head against wall", app.run)
	if err != nil {
		return err
	}

	// Optional
	cli.SetDescription(`
Long description.
Could be multi-line.`)

	// Optional
	cli.SetExamples(`
One or more examples.

Could be multi-line.`)

	// Optional
	cli.SetFooter("For more information, see https://www.example.org/")

	if err := cli.AddFlags(
		&clim.Flag{
			Value: clim.Int(&app.count, 3),
			Short: "c", Long: "count", Label: "N", Help: "How many times",
		},
		&clim.Flag{
			Value: clim.String(&app.wall, "cardboard"),
			Long:  "wall", Help: "Type of wall",
		},
		&clim.Flag{
			Value: clim.Bool(&app.dryRun, false),
			Long:  "dry-run", Help: "Enable dry-run",
		},
		&clim.Flag{
			Value: clim.IntSlice(&app.windows, nil),
			Long:  "windows", Label: "N[,N,..]",
			Help: "Windows sequence",
		},
		&clim.Flag{
			Value: clim.IntSlice(&app.doors, nil),
			Long:  "doors", Label: "N[,N,..]",
			Help: "Doors sequence",
		},
		&clim.Flag{
			Value: clim.StringSlice(&app.floors, nil),
			Long:  "floors", Label: "F[,F,..]",
			Help: "Floors sequence",
		},
	); err != nil {
		return err
	}

	action, err := cli.Parse(args)
	if err != nil {
		return err
	}

	return action(0)
}

func (args *Application) run(uctx int) error {
	// Validation
	if clim.CountTrue(args.doors != nil, args.windows != nil,
		args.floors != nil) > 1 {
		return clim.NewParseError("only one of doors, windows, floors can be specified")
	}

	for i := range args.count {
		fmt.Println(i+1, "flatten against", args.wall)
	}
	return nil
}
