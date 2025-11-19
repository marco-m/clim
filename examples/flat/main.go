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
	cli, err := clim.NewTop("flat", "flattens head against wall")
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

	if _, err := cli.Parse(args); err != nil {
		return err
	}

	return app.run()
}

func (app *Application) run() error {
	// Validation
	if clim.CountTrue(app.doors != nil, app.windows != nil,
		app.floors != nil) > 1 {
		return clim.NewParseError("only one of doors, windows, floors can be specified")
	}

	for i := range app.count {
		fmt.Println(i+1, "flatten against", app.wall)
	}
	return nil
}
