package clim_test

import (
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestSimpleHelp(t *testing.T) {
	type Args struct {
		count  int
		wall   string
		dryRun bool
	}
	var args Args
	cli, err := clim.New[any](nil, "bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{
			Value: clim.Int(&args.count, 3),
			Short: "c",
			Long:  "count",
			Label: "N",
			Help:  "How many times",
		},
		&clim.Flag{
			Value: clim.String(&args.wall, "cardboard"),
			// Short is optional, here we don't set it.
			Long: "wall",
			// Default for Label: uppercase(Long)
			Help: "Type of wall",
		},
		&clim.Flag{
			Value: clim.Bool(&args.dryRun, false),
			Long:  "dry-run",
			Help:  "Enable dry-run",
		})
	rosina.AssertNoError(t, err)

	want := `bang -- bangs head against wall

Usage: bang [options]

Options:

 -c, --count N    How many times (default: 3)
 --dry-run        Enable dry-run (default: false)
 --wall WALL      Type of wall (default: cardboard)

 -h, --help       Print this help and exit
`

	_, err = cli.Parse([]string{"-h"})

	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertErrorContains(t, err, want)
}

func TestHelpOfRequiredFlag(t *testing.T) {
	var count int
	var level int

	cli, err := clim.New[any](nil, "bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{
			// Default value with Required, will be ignored also in the help.
			Value:    clim.Int(&count, 3),
			Long:     "count",
			Required: true,
		},
		&clim.Flag{
			// Default value without Required, normal handling.
			Value: clim.Int(&level, 5),
			Long:  "level",
		})
	rosina.AssertNoError(t, err)

	want := `bang -- bang head

Usage: bang [options]

Options:

 --count COUNT     (required)
 --level LEVEL     (default: 5)

 -h, --help       Print this help and exit
`

	_, err = cli.Parse([]string{"-h"})

	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertErrorContains(t, err, want)
}

func TestHelpOfOptionalFields(t *testing.T) {
	want := `bang -- bang head

 this is the description

Usage: bang [options]

Examples:

 This is a multi-line example.

 This is the last line of the example.

Options:

 -h, --help    Print this help and exit

 this is the footer
`

	cli, err := clim.New[any](nil, "bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	cli.SetDescription("this is the description")
	cli.SetExamples(`
This is a multi-line example.

This is the last line of the example.`)
	cli.SetFooter("this is the footer")

	_, err = cli.Parse([]string{"-h"})

	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertTextEqual(t, err.Error(), want, "help message")
}

func TestHelpSubCommandsOneLevel(t *testing.T) {
	cli, err := clim.New[any](nil, "bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.New[any](cli, "sub", "I am a subcommand", nil)
	rosina.AssertNoError(t, err)

	want := `bang -- bangs head against wall

Usage: bang <command> [options]

Commands:

 sub     I am a subcommand

Options:

 -h, --help    Print this help and exit
`
	_, err = cli.Parse([]string{"-h"})

	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertDeepEqual(t, err.Error(), want, "error text")
}

func TestHelpSubCommandsTwoLevels(t *testing.T) {
	cli, err := clim.New[any](nil, "bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	sub1, err := clim.New[any](cli, "sub1", "I am a subcommand at level 1", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.New[any](sub1, "sub2", "I am a subcommand at level 2", nil)
	rosina.AssertNoError(t, err)

	want := `bang sub1 sub2 -- I am a subcommand at level 2

Usage: bang sub1 sub2 [options]

Options:

 -h, --help    Print this help and exit
`
	_, err = cli.Parse([]string{"sub1", "sub2", "-h"})

	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertDeepEqual(t, err.Error(), want, "error text")
}

func TestHelpSubCommandsGroup(t *testing.T) {
	cli, err := clim.New[any](nil, "bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	subCliA, err := clim.New[any](cli, "sub-A", "I am subcommand A", nil)
	rosina.AssertNoError(t, err)

	subCliB, err := clim.New[any](cli, "sub-B", "I am subcommand B", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddGroup("group 1", subCliA)
	rosina.AssertNoError(t, err)

	err = cli.AddGroup("group 2", subCliB)
	rosina.AssertNoError(t, err)

	want := `bang -- bangs head against wall

Usage: bang <command> [options]

available commands:

group 1:

 sub-A     I am subcommand A

group 2:

 sub-B     I am subcommand B

Options:

 -h, --help    Print this help and exit
`
	_, err = cli.Parse([]string{"-h"})

	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertDeepEqual(t, err.Error(), want, "error text")
}

func TestPosArgsRequiredHelpTake1(t *testing.T) {
	want := `bang -- bang head

Usage: bang [options] COUNT NAME COLOR...

Options:

 -h, --help    Print this help and exit

Positional arguments:

 COUNT         How many foos (required)
 NAME          Name of the foos (required)
 COLOR...      One or more colors (required)
`

	cli, err := clim.New[any](nil, "bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	var positionals []string
	err = cli.AddPosArgs(&positionals,
		clim.Pair{"COUNT", "How many foos (required)"},
		clim.Pair{"NAME", "Name of the foos (required)"},
		clim.Pair{"COLOR...", "One or more colors (required)"})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"-h"})
	rosina.AssertErrorIs(t, err, clim.ErrHelp)
	rosina.AssertDeepEqual(t, err.Error(), want, "help message")
}

// func TestPosArgsSimpleHelp(t *testing.T) {
// 	want := `bang -- bang head

// Usage: bang [options] NAME COUNT

// Options:

//  -h, --help    Print this help and exit

// Positional arguments:

//  NAME       Name of the foos
//  COUNT      How many foos
// `
// 	var foos string
// 	cli := clim.New[any]("bang", "bang head", nil)
// 	err := cli.AddPosArgs(&clim.PosArg{
// 		Value: clim.String(&foos, ""),
// 		Name:  "NAME", Help: "Name of the foos", Required: true,
// 	})
// 	rosina.AssertNoError(t, err)

// 	_, err = cli.Parse([]string{"-h"})

// 	rosina.AssertErrorIs(t, err, clim.ErrHelp)
// 	rosina.AssertDeepEqual(t, err.Error(), want, "help message")
// }
