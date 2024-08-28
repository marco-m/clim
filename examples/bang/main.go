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

type Application struct {
	count    int
	wall     string
	dryRun   bool
	sequence []int
}

func mainErr(args []string) error {
	var app Application
	cli := clim.New("bang", "bangs head against wall", app.run)

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

	cli.AddFlag(&clim.Flag{
		Value: clim.Int(&app.count, 3),
		Short: "c", Long: "count", Label: "N", Help: "How many times",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.String(&app.wall, "cardboard"),
		Long:  "wall", Help: "Type of wall",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&app.dryRun, false),
		Long:  "dry-run", Help: "Enable dry-run",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.IntSlice(&app.sequence, []int{1, 2, 3}),
		Short: "s", Long: "sequence", Label: "N[,N,..]",
		Help: "bang sequence",
	})

	action, err := cli.Parse(args)
	if err != nil {
		return err
	}

	return action(0)
}

func (args *Application) run(uctx int) error {
	for i := range args.count {
		fmt.Println(i+1, "bang against", args.wall)
	}
	return nil
}
