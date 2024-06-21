package clim_test

import (
	"strings"
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
	cli := clim.New("bang", "bangs head against wall")
	cli.IntVar(&args.count, "c", "count", "N", 3, "How many times")
	cli.StringVar(&args.wall, "", "wall", "WALL", "cardboard", "Type of wall")
	cli.BoolVar(&args.dryRun, "", "dry-run", false, "Enable dry-run")

	want := strings.TrimSpace(`
bang -- bangs head against wall

Usage: bang [options]

Options:

 -c, --count N    How many times (default: 3)
 --dry-run        Enable dry-run (default: false)
 --wall WALL      Type of wall (default: cardboard)

 -h, --help       Print this help and exit
`)

	_, err := cli.Parse([]string{"-h"})

	qt.Assert(t, qt.ErrorIs(err, clim.ErrHelp))
	qt.Assert(t, qt.Equals(err.Error(), want))
}

func TestDestinationsMustBeUnique(t *testing.T) {
	var count int
	cli := clim.New("banana", "I am tasty")

	// 1st reference to variable 'count': OK.
	cli.IntVar(&count, "", "count", "COUNT", 3, "How many")

	// 2nd reference to variable 'count': panic
	qt.Assert(t, qt.PanicMatches(func() {
		cli.IntVar(&count, "x", "extra", "EXTRA", 0, "Extra taste")
	}, `long flag name "extra": variable already bound to flag "count"`))
}

func TestLongFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli := clim.New("banana", "I am tasty")

	// 1st long flag '--count'
	cli.IntVar(&count, "", "count", "COUNT", 3, "How many")

	// 2nd long flag '--count' panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.IntVar(&extra, "", "count", "EXTRA", 0, "Extra taste")
	}, `banana: long flag name "count" already defined`))
}

func TestShortFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli := clim.New("banana", "I am tasty")

	// 1st short flag '-c'
	cli.IntVar(&count, "c", "count", "COUNT", 3, "How many")

	// 2nd short flag '-c' panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.IntVar(&extra, "c", "extra", "EXTRA", 0, "Extra taste")
	}, `banana: short flag name "c" already defined`))
}

func TestCannotOverrideHelpFlag(t *testing.T) {
	var count int
	var extra int
	cli := clim.New("banana", "I am tasty")

	// Attempt to override '--help' panics
	// FIXME In the future I would like to allow to ovverride --help
	//       to allow the program to provide more verbose information?
	qt.Assert(t, qt.PanicMatches(func() {
		cli.IntVar(&count, "x", "help", "COUNT", 0, "How many")
	}, `cannot override long flag name "help"`))

	// Attempt to override '-h' panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.IntVar(&extra, "h", "extra", "EXTRA", 0, "Extra taste")
	}, `cannot override short flag name "h"`))
}

func TestLongFlagIsMandatory(t *testing.T) {
	var count int
	var extra int
	cli := clim.New("banana", "I am tasty")

	// Empty short flag is OK
	cli.IntVar(&count, "", "count", "COUNT", 3, "How many")

	// Empty long flag panics
	qt.Assert(t, qt.PanicMatches(func() {
		cli.IntVar(&extra, "x", "", "EXTRA", 4, "Extra taste")
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
		cli := clim.New("banana", "I am tasty")
		var count int

		qt.Assert(t, qt.PanicMatches(func() {
			cli.IntVar(&count, tc.short, tc.long, "N", 3, "How many")
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

func TestNameCannotBeEmpty(t *testing.T) {
	qt.Assert(t, qt.PanicMatches(
		func() {
			clim.New("", "I am tasty")
		}, `clim\.New: name cannot be empty`,
	))
}

func TestParseOnePairSuccess(t *testing.T) {
	var count int
	cli := clim.New("basket", "juicy fruits")
	cli.IntVar(&count, "", "count", "N", 3, "How many")

	_, err := cli.Parse([]string{"--count", "42"})

	qt.Check(t, qt.IsNil(err))
	qt.Check(t, qt.Equals(count, 42))
}

func TestParseOnePairUnrecognized(t *testing.T) {
	var count int
	cli := clim.New("basket", "juicy fruits")
	cli.IntVar(&count, "", "count", "N", 3, "How many")

	_, err := cli.Parse([]string{"--fruit", "42"})

	qt.Check(t, qt.ErrorMatches(err, `unrecognized flag "--fruit"`))
}

func TestArgs(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want []string
	}

	test := func(t *testing.T, tc testCase) {
		var count int
		cli := clim.New("basket", "juicy fruits")
		cli.IntVar(&count, "", "count", "N", 3, "How many")

		_, err := cli.Parse(tc.args)

		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.DeepEquals(cli.Args(), tc.want))
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
