package clim_test

import (
	"errors"
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestVariableCanBeBoundOnlyOnce(t *testing.T) {
	var count int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// 1st reference to variable 'count': OK.
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	// 2nd reference to variable 'count': panic
	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 0), Long: "extra"})
	},
		`long flag name "extra": variable already bound to flag "count"`)
}

func TestLongFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// 1st long flag '--count'
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	// 2nd long flag '--count' panics
	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&extra, 0), Long: "count"})
	},
		`banana: long flag name "count" already defined`)
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
	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{
			Value: clim.Int(&extra, 0),
			Short: "c", Long: "extra",
		})
	},
		`banana: short flag name "c" already defined`)
}

func TestShortFlagMustBeOneChar(t *testing.T) {
	var count int
	cli := clim.New[any]("banana", "I am tasty", nil)

	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 0), Short: "co", Long: "help"})
	},
		`short flag name "co" must be exactly 1 character`)
}

func TestLongFlagMustBeMoreThanOneChar(t *testing.T) {
	var count int
	cli := clim.New[any]("banana", "I am tasty", nil)

	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 0), Long: "c"})
	},
		`long flag name "c" must be at least 2 characters`)
}

func TestCannotOverrideHelpFlag(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// Attempt to override '--help' panics
	// FIXME In the future I would like to allow to ovverride --help
	//       to allow the program to provide more verbose information?
	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 0), Long: "help"})
	},
		`cannot override long flag name "help"`)

	// Attempt to override '-h' panics
	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{
			Value: clim.Int(&extra, 0),
			Short: "h", Long: "extra",
		})
	},
		`cannot override short flag name "h"`)
}

func TestLongFlagIsMandatory(t *testing.T) {
	var count int
	var extra int
	cli := clim.New[any]("banana", "I am tasty", nil)

	// Empty short flag is OK
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	// Empty long flag panics
	rosina.AssertPanicTextContains(t, func() {
		cli.AddFlag(&clim.Flag{Value: clim.Int(&extra, 4), Short: "x"})
	},
		`long flag name cannot be empty`)
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

		rosina.AssertPanicTextContains(t, func() {
			cli.AddFlag(&clim.Flag{
				Value: clim.Int(&count, 3),
				Short: tc.short, Long: tc.long,
			})
		},
			tc.want)
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
	rosina.AssertPanicTextContains(t, func() {
		clim.New[any]("", "I am tasty", nil)
	},
		`clim.New: name cannot be empty`)
}

func TestActionMissing(t *testing.T) {
	cli := clim.New[string]("basket", "juicy fruits", nil)
	action, err := cli.Parse(nil)

	rosina.AssertIsNil(t, err)
	err = action("hello")
	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, "basket: no action registered")
}

func TestActionPresent(t *testing.T) {
	cli := clim.New[string]("basket", "juicy fruits",
		func(uctx string) error { return errors.New(uctx) })

	// In this simple case, it might be unclear why the indirection
	// of passing through action. It becomes evident when using subcommands.
	action, err := cli.Parse(nil)
	rosina.AssertIsNil(t, err)

	err = action("mango")
	rosina.AssertIsNotNil(t, err)
	rosina.AssertTextEqual(t, err.Error(), "mango", "action error message")
}

func TestParseOneFlagPairSuccess(t *testing.T) {
	var count int
	cli := clim.New[any]("basket", "juicy fruits", nil)
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	_, err := cli.Parse([]string{"--count", "42"})

	rosina.AssertIsNil(t, err)
	rosina.AssertEqual(t, count, 42, "count")
}

func TestParseOneFlagPairUnrecognized(t *testing.T) {
	var count int
	cli := clim.New[any]("basket", "juicy fruits", nil)
	cli.AddFlag(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})

	_, err := cli.Parse([]string{"--fruit", "42"})
	rosina.AssertErrorContains(t, err, `unrecognized flag "--fruit"`)

	_, err = cli.Parse([]string{"--fruit"})
	rosina.AssertErrorContains(t, err, `unrecognized flag "--fruit"`)

	_, err = cli.Parse([]string{"-f", "42"})
	rosina.AssertErrorContains(t, err, `unrecognized flag "-f"`)

	_, err = cli.Parse([]string{"-f"})
	rosina.AssertErrorContains(t, err, `unrecognized flag "-f"`)
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

	rosina.AssertIsNil(t, err)
	rosina.AssertEqual(t, count, 1, "count (parsed)")
	rosina.AssertEqual(t, level, 5, "level (default value)")
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

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `missing required options: count, foo`)
}

func TestSubCommandWithRequiredOptionFailure(t *testing.T) {
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

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `missing required options: foo`)
}

func TestMissingSubcommandFailure(t *testing.T) {
	cli := clim.New[any]("bang", "bang head", nil)
	subCli := clim.New[any]("sub", "I am a subcommand", nil)
	cli.AddCLI(subCli)

	_, err := cli.Parse([]string{})

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `expected a command`)
}

func TestWrongSubcommandFailure(t *testing.T) {
	cli := clim.New[any]("bang", "bang head", nil)
	subCli := clim.New[any]("sub", "I am a subcommand", nil)
	cli.AddCLI(subCli)

	_, err := cli.Parse([]string{"hello"})

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `unrecognized command "hello"`)
}

func TestSubCommandNamesMustBeUnique(t *testing.T) {
	cli := clim.New[any]("bang", "bang head", nil)
	subCliA := clim.New[any]("sub", "I am a subcommand A", nil)
	cli.AddCLI(subCliA)
	subCliB := clim.New[any]("sub", "I am a subcommand B", nil)

	rosina.AssertPanicTextContains(t, func() {
		cli.AddCLI(subCliB)
	},
		`bang: subcommand "sub" already defined`)
}

func TestAddGroupSuccess(t *testing.T) {
	child := clim.New[any]("child", "I am a child", nil)
	root := clim.New[any]("bang", "bang head", nil)
	root.AddGroup("ciccio", root.AddCLI(child))
}

func TestAddGroupMissingChildren(t *testing.T) {
	root := clim.New[any]("bang", "bang head", nil)
	rosina.AssertPanicTextContains(t, func() {
		root.AddGroup("ciccio")
	},
		`AddGroup ciccio: child list is empty`)
}

func TestAddGroupMissingAddCLI(t *testing.T) {
	// AddGroup ciccio: child child is missing previous AddCLI
	child := clim.New[any]("child", "I am a child", nil)
	root := clim.New[any]("bang", "bang head", nil)
	rosina.AssertPanicTextContains(t, func() {
		root.AddGroup("ciccio", (child))
	},
		`AddGroup ciccio: child child is missing previous AddCLI`)
}

func TestContTrue(t *testing.T) {
	type testCase struct {
		name string
		args []bool
		want int
	}

	test := func(t *testing.T, tc testCase) {
		have := clim.CountTrue(tc.args...)
		rosina.AssertEqual(t, have, tc.want, "CountTrue")
	}

	testCases := []testCase{
		{name: "empty", args: []bool{}, want: 0},
		{name: "1 true, 1 false", args: []bool{true, false}, want: 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}
