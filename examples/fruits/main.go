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
	err := mainErr()
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

type Args struct {
	count  int
	wall   string
	dryRun bool
}

func mainErr() error {
	var args Args
	cli := clim.New("bang", "bangs head against wall")
	cli.Action(func() error { return run(args) })

	cli.AddFlag(clim.IntVal(&args.count, 3),
		"c", "count", "N", "How many times")
	cli.AddFlag(clim.StringVal(&args.wall, "cardboard"),
		"", "wall", "WALL", "Type of wall")
	cli.AddFlag(clim.BoolVal(&args.dryRun, false),
		"", "dry-run", "", "Enable dry-run")

	action, err := cli.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return action()
}

func run(args Args) error {
	for i := range args.count {
		fmt.Println(i+1, "bang against", args.wall)
	}
	return nil
}
