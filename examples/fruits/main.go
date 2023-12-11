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

func mainErr() error {
	type Args struct {
		count  int
		wall   string
		dryRun bool
	}
	var args Args
	cli := clim.New("bang", "bangs head against wall")
	cli.IntVar(&args.count, "c", "count", "N", 3, "How many times")
	cli.StringVar(&args.wall, "", "wall", "WALL", "cardboard", "Type of wall")
	cli.BoolVar(&args.dryRun, "", "dry-run", false, "Enable dry-run")

	if err := cli.Parse(os.Args[1:]); err != nil {
		return err
	}

	for i := range args.count {
		fmt.Println(i+1, "bang against", args.wall)
	}
	return nil
}
