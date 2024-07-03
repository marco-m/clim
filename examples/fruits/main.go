package main

import (
	"context"
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
	count    int
	wall     string
	dryRun   bool
	sequence []int
}

func mainErr() error {
	var args Args
	cli := clim.New("bang", "bangs head against wall")
	cli.SetAction(args.run)

	cli.AddFlag(&clim.Flag{Value: clim.Int(&args.count, 3),
		Short: "c", Long: "count", Label: "N", Desc: "How many times"})
	cli.AddFlag(&clim.Flag{Value: clim.String(&args.wall, "cardboard"),
		Long: "wall", Desc: "Type of wall"})
	cli.AddFlag(&clim.Flag{Value: clim.Bool(&args.dryRun, false),
		Long: "dry-run", Desc: "Enable dry-run"})
	cli.AddFlag(&clim.Flag{Value: clim.IntSlice(&args.sequence, []int{1, 2, 3}),
		Short: "s", Long: "sequence", Label: "N[,N,..]",
		Desc: "bang sequence"})

	action, err := cli.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	ctx := context.Background()

	return action(ctx)
}

func (args *Args) run(ctx context.Context) error {
	for i := range args.count {
		fmt.Println(i+1, "bang against", args.wall)
	}
	return nil
}
