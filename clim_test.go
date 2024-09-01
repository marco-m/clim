package clim_test

import (
	"testing"

	"github.com/go-quicktest/qt"

	"github.com/marco-m/clim"
)

func TestSimpleHelp(t *testing.T) {
	type Args struct {
		count  int
		wall   string
		dryRun bool
	}
	var args Args
	cli := clim.New[any]("bang", "bangs head against wall", nil)

	cli.AddFlag(&clim.Flag{
		Value: clim.Int(&args.count, 3),
		Short: "c",
		Long:  "count",
		Label: "N",
		Help:  "How many times",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.String(&args.wall, "cardboard"),
		// Short is optional, here we don't set it.
		Long: "wall",
		// Default for Label: uppercase(Long)
		Help: "Type of wall",
	})
	cli.AddFlag(&clim.Flag{
		Value: clim.Bool(&args.dryRun, false),
		Long:  "dry-run",
		Help:  "Enable dry-run",
	})

	want := `bang -- bangs head against wall

Usage: bang [options]

Options:

 -c, --count N    How many times (default: 3)
 --dry-run        Enable dry-run (default: false)
 --wall WALL      Type of wall (default: cardboard)

 -h, --help       Print this help and exit
`

	_, err := cli.Parse([]string{"-h"})

	qt.Assert(t, qt.ErrorIs(err, clim.ErrHelp))
	qt.Assert(t, qt.Equals(err.Error(), want))
}

func TestVariableCanBeBoundOnlyOnce(t *testing.T) {
	var count int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// 1st reference to variable 'count': OK.
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	// 2nd reference to variable 'count': panic
	qt.Assert(t, qt.PanicMatches(func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 0), Long: "extra"})
	}, `long flag name "extra": variable already bound to flag "count"`))
}

func TestLongFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// 1st long flag '--count'
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	// 2nd long flag '--count' panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&extra, 0), Long: "count"})
	}, `banana: long flag name "count" already defined`))
}

func TestShortFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// 1st short flag '-c'
	cli.AddFlag(&clim.Flag{
		Value: clim.Int(&count, 3),
		Short: "c", Long: "count",
	})

	// 2nd short flag '-c' panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.AddFlag(&clim.Flag{
			Value: clim.Int(&extra, 0),
			Short: "c", Long: "extra",
		})
	}, `banana: short flag name "c" already defined`))
}

func TestCannotOverrideHelpFlag(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// Attempt to override '--help' panics
	// FIXME In the future I would like to allow to ovverride --help
	//       to allow the program to provide more verbose information?
	qt.Assert(t, qt.PanicMatches(func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 0), Long: "help"})
	}, `cannot override long flag name "help"`))

	// Attempt to override '-h' panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.AddFlag(&clim.Flag{
			Value: clim.Int(&extra, 0),
			Short: "h", Long: "extra",
		})
	}, `cannot override short flag name "h"`))
}

func TestLongFlagIsMandatory(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// Empty short flag is OK
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	// Empty long flag panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&extra, 4), Short: "x"})
	}, `long flag name cannot be empty`))
}

func TestFlagsNamingConstraints(t *testing.T) {
	type testCase struct {
		name  string
		short string
		long  string
		want  string
	}

	test := func(t *testing.T, tc testCase) {
		cli := clim.New[any]("banana", "I am tasty", nil)
		var count int

		qt.Assert(t, qt.PanicMatches(func() {
			cli.AddFlag(&clim.Flag{
				Value: clim.Int(&count, 3),
				Short: tc.short, Long: tc.long,
			})
		}, tc.want))
	}

	testCases := []testCase{
		{
			name: "long begins with '-'",
			long: "-foo",
			want: `long flag name "-foo" must not begin with '-'`,
		},
		{
			name:  "short begins with '-'",
			short: "-",
			long:  "long",
			want:  `short flag name must not begin with '-'`,
		},
		{
			name: "long contains '='",
			long: "foo=bar",
			want: `long flag name "foo=bar" must not contain '='`,
		},
		{
			name:  "short contains '='",
			short: "=",
			long:  "foobar",
			want:  `short flag name must not contain '='`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestCliNameCannotBeEmpty(t *testing.T) {
	qt.Assert(t, qt.PanicMatches(
		func() {
			clim.New[any]("", "I am tasty", nil)
		}, `clim\.New: name cannot be empty`,
	))
}

func TestParseOneFlagPairSuccess(t *testing.T) {
	var count int
	cli := clim.New[any]("basket", "juicy fruits", nil)
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	_, err := cli.Parse([]string{"--count", "42"})

	qt.Check(t, qt.IsNil(err))
	qt.Check(t, qt.Equals(count, 42))
}

func TestParseOneFlagPairUnrecognized(t *testing.T) {
	var count int
	cli := clim.New[any]("basket", "juicy fruits", nil)
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	_, err := cli.Parse([]string{"--fruit", "42"})

	qt.Check(t, qt.ErrorMatches(err, `unrecognized flag "--fruit"`))
}

func TestPosArgs(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want []string
	}

	test := func(t *testing.T, tc testCase) {
		var count int
		cli := clim.New[any]("basket", "juicy fruits", nil)
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

		_, err := cli.Parse(tc.args)

		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.DeepEquals(cli.PosArgs(), tc.want))
	}

	testCases := []testCase{
		{
			name: "no positionals",
			args: []string{"--count=42"},
			want: []string{},
		},
		{
			name: "vanilla",
			args: []string{"--count=42", "a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "sneaky",
			args: []string{"--count=42", "---a", "b"},
			want: []string{"---a", "b"},
		},
		{
			name: "after the first positional, a flag is not a flag",
			args: []string{"--count=42", "a", "-b"},
			want: []string{"a", "-b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestRequiredHelp(t *testing.T) {
	var count int
	var level int
	cli := clim.New[any]("bang", "bang head", nil)
	cli.AddFlag(&clim.Flag{
		// Default value with Required, will be ignored also in the help.
		Value:    clim.Int(&count, 3),
		Long:     "count",
		Required: true,
	})
	cli.AddFlag(&clim.Flag{
		// Default value without Required, normal handling.
		Value: clim.Int(&level, 5),
		Long:  "level",
	})

	want := `bang -- bang head

Usage: bang [options]

Options:

 --count COUNT     (required)
 --level LEVEL     (default: 5)

 -h, --help       Print this help and exit
`

	_, err := cli.Parse([]string{"-h"})

	qt.Assert(t, qt.ErrorIs(err, clim.ErrHelp))
	qt.Assert(t, qt.Equals(err.Error(), want))
}

func TestRequiredIgnoresDefaultSuccess(t *testing.T) {
	var count int
	var level int
	cli := clim.New[any]("bang", "bang head", nil)
	cli.AddFlag(&clim.Flag{
		// Default value with Required, will be ignored also in the help.
		Value:    clim.Int(&count, 3),
		Long:     "count",
		Required: true,
	})
	cli.AddFlag(&clim.Flag{
		// Default value without Required, normal handling.
		Value: clim.Int(&level, 5),
		Long:  "level",
	})

	_, err := cli.Parse([]string{"--count=1"})

	qt.Check(t, qt.IsNil(err))
	qt.Check(t, qt.Equals(count, 1)) // parsed
	qt.Check(t, qt.Equals(level, 5)) // from default value
}

func TestRequiredFailure(t *testing.T) {
	var count int
	var level int
	var foo int
	cli := clim.New[any]("bang", "bang head", nil)
	cli.AddFlag(&clim.Flag{
		Value:    clim.Int(&count, 3),
		Long:     "count",
		Required: true,
	})
	cli.AddFlag(&clim.Flag{
		// Default value without Required, normal handling.
		Value: clim.Int(&level, 5),
		Long:  "level",
	})
	cli.AddFlag(&clim.Flag{
		Value:    clim.Int(&foo, 3),
		Long:     "foo",
		Required: true,
	})

	_, err := cli.Parse(nil)

	qt.Assert(t, qt.ErrorIs(err, clim.ErrParse))
	qt.Assert(t, qt.Equals(err.Error(), `missing required options: count, foo`))
}

func TestSubCommandRequiredFailure(t *testing.T) {
	var count int
	var foo int

	cli := clim.New[any]("bang", "bang head", nil)
	cli.AddFlag(&clim.Flag{
		Value:    clim.Int(&count, 0),
		Long:     "count",
		Required: true,
	})

	subCli := clim.New[any]("sub", "I am a subcommand", nil)
	subCli.AddFlag(&clim.Flag{
		Value:    clim.Int(&foo, 0),
		Long:     "foo",
		Required: true,
	})
	cli.AddCLI(subCli)

	_, err := cli.Parse([]string{"--count=22", "sub"})

	qt.Assert(t, qt.ErrorIs(err, clim.ErrParse))
	qt.Assert(t, qt.Equals(err.Error(), `missing required options: foo`))
}
