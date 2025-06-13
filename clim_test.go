package clim_test

import (
	"errors"
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestVariableCanBeBoundOnlyOnce(t *testing.T) {
	var count int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{Value: clim.Int(&count, 3), Long: "count"},
		&clim.Flag{Value: clim.Int(&count, 0), Long: "extra"},
	)
	rosina.AssertErrorContains(t, err,
		`long flag name "extra": variable already bound to flag "count"`)

	rosina.AssertErrorContains(t, err,
		`long flag name "extra": variable already bound to flag "count"`)
}

func TestLongFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{Value: clim.Int(&count, 3), Long: "count"},
		&clim.Flag{Value: clim.Int(&extra, 0), Long: "count"},
	)
	rosina.AssertErrorContains(t, err,
		`banana: long flag name "count" already defined`)
}

func TestShortFlagsMustBeUnique(t *testing.T) {
	var count int
	var extra int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{
			Value: clim.Int(&count, 3),
			Short: "c", Long: "count",
		},
		&clim.Flag{
			Value: clim.Int(&extra, 0),
			Short: "c", Long: "extra",
		})

	rosina.AssertErrorContains(t, err,
		`banana: short flag name "c" already defined`)
}

func TestShortFlagMustBeOneChar(t *testing.T) {
	var count int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{Value: clim.Int(&count, 0), Short: "co", Long: "help"},
	)

	rosina.AssertErrorContains(t, err,
		`short flag name "co" must be exactly 1 character`)
}

func TestLongFlagMustBeMoreThanOneChar(t *testing.T) {
	var count int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{Value: clim.Int(&count, 0), Long: "c"})

	rosina.AssertErrorContains(t, err,
		`long flag name "c" must be at least 2 characters`)
}

func TestCannotOverrideLongHelpFlag(t *testing.T) {
	var count int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	// FIXME In the future I would like to allow to ovverride --help
	//       to allow the program to provide more verbose information?
	err = cli.AddFlags(&clim.Flag{Value: clim.Int(&count, 0), Long: "help"})

	rosina.AssertErrorContains(t, err,
		`cannot override long flag name "help"`)
}

func TestCannotOverrideShortHelpFlag(t *testing.T) {
	var extra int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value: clim.Int(&extra, 0),
		Short: "h", Long: "extra",
	})

	rosina.AssertErrorContains(t, err,
		`cannot override short flag name "h"`)
}

func TestLongFlagIsMandatory(t *testing.T) {
	var count int
	var extra int
	cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{Value: clim.Int(&count, 3), Long: "count"},
		&clim.Flag{Value: clim.Int(&extra, 4), Short: "x"})

	rosina.AssertErrorContains(t, err, `long flag name cannot be empty`)
}

func TestFlagsNamingConstraints(t *testing.T) {
	type testCase struct {
		name  string
		short string
		long  string
		want  string
	}

	test := func(t *testing.T, tc testCase) {
		cli, err := clim.NewTop[any]("banana", "I am tasty", nil)
		rosina.AssertNoError(t, err)

		var count int
		err = cli.AddFlags(&clim.Flag{
			Value: clim.Int(&count, 3),
			Short: tc.short, Long: tc.long,
		})
		rosina.AssertErrorContains(t, err, tc.want)
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
	_, err := clim.NewTop[any]("", "I am tasty", nil)
	rosina.AssertErrorContains(t, err, `cli name cannot be empty`)
}

func TestActionMissing(t *testing.T) {
	cli, err := clim.NewTop[string]("basket", "juicy fruits", nil)
	rosina.AssertNoError(t, err)

	action, err := cli.Parse(nil)

	rosina.AssertNoError(t, err)
	err = action("hello")
	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, "basket: no action registered")
}

func TestActionPresent(t *testing.T) {
	cli, err := clim.NewTop[string]("basket", "juicy fruits",
		func(uctx string) error { return errors.New(uctx) })
	rosina.AssertNoError(t, err)

	// In this simple case, it might be unclear why the indirection
	// of passing through action. It becomes evident when using subcommands.
	action, err := cli.Parse(nil)
	rosina.AssertNoError(t, err)

	err = action("mango")
	rosina.AssertErrorContains(t, err, "mango")
}

func TestParseOneFlagPairSuccess(t *testing.T) {
	var count int
	cli, err := clim.NewTop[any]("basket", "juicy fruits", nil)
	rosina.AssertNoError(t, err)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--count", "42"})

	rosina.AssertNoError(t, err)
	rosina.AssertEqual(t, count, 42, "count")
}

func TestParseOneFlagPairUnrecognized(t *testing.T) {
	var count int
	cli, err := clim.NewTop[any]("basket", "juicy fruits", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{Value: clim.Int(&count, 3), Long: "count"})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--fruit", "42"})
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
	cli, err := clim.NewTop[any]("bang", "bang head", nil)
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

	_, err = cli.Parse([]string{"--count=1"})

	rosina.AssertNoError(t, err)
	rosina.AssertEqual(t, count, 1, "count (parsed)")
	rosina.AssertEqual(t, level, 5, "level (default value)")
}

func TestRequiredFailure(t *testing.T) {
	var count int
	var level int
	var foo int
	cli, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(
		&clim.Flag{
			Value:    clim.Int(&count, 3),
			Long:     "count",
			Required: true,
		},
		&clim.Flag{
			// Default value without Required, normal handling.
			Value: clim.Int(&level, 5),
			Long:  "level",
		},
		&clim.Flag{
			Value:    clim.Int(&foo, 3),
			Long:     "foo",
			Required: true,
		})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse(nil)

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `missing required options: count, foo`)
}

func TestSubCommandWithRequiredOptionFailure(t *testing.T) {
	var count int
	var foo int

	cli, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value:    clim.Int(&count, 0),
		Long:     "count",
		Required: true,
	})
	rosina.AssertNoError(t, err)

	subCli, err := clim.NewSub[any](cli, "sub", "I am a subcommand", nil)
	rosina.AssertNoError(t, err)

	err = subCli.AddFlags(&clim.Flag{
		Value:    clim.Int(&foo, 0),
		Long:     "foo",
		Required: true,
	})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--count=22", "sub"})

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `missing required options: foo`)
}

func TestMissingSubcommandFailure(t *testing.T) {
	cli, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.NewSub[any](cli, "sub", "I am a subcommand", nil)
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{})

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `expected a command`)
}

func TestWrongSubcommandFailure(t *testing.T) {
	cli, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.NewSub[any](cli, "sub", "I am a subcommand", nil)
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"hello"})

	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err, `unrecognized command "hello"`)
}

func TestSubCommandNamesMustBeUnique(t *testing.T) {
	cli, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.NewSub[any](cli, "sub", "I am a subcommand A", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.NewSub[any](cli, "sub", "I am a subcommand B", nil)
	rosina.AssertErrorContains(t, err, `bang: subcommand "sub" already defined`)
}

func TestCannotAddSubCommandAfterPosArgs(t *testing.T) {
	cli, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	var positionals []string
	err = cli.AddPosArgs(&positionals, clim.Pair{"NAME", "Name of the foos"})
	rosina.AssertNoError(t, err)

	_, err = clim.NewSub[any](cli, "sub", "I am a subcommand A", nil)
	rosina.AssertErrorContains(t, err,
		`bang: already have pos args; cannot have also subcommand "sub"`)
}

func TestAddGroupSuccess(t *testing.T) {
	root, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	child, err := clim.NewSub[any](root, "child", "I am a child", nil)
	rosina.AssertNoError(t, err)

	err = root.AddGroup("ciccio", child)
	rosina.AssertNoError(t, err)
}

func TestAddGroupMissingChildren(t *testing.T) {
	root, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	err = root.AddGroup("ciccio")
	rosina.AssertErrorContains(t, err, `AddGroup ciccio: child list is empty`)
}

func TestAddGroupMissingAddCLI(t *testing.T) {
	// AddGroup ciccio: child child is missing previous AddCLI
	child, err := clim.NewTop[any]("child", "I am a child", nil)
	rosina.AssertNoError(t, err)

	root, err := clim.NewTop[any]("bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	err = root.AddGroup("ciccio", child)
	rosina.AssertErrorContains(t, err,
		"AddGroup ciccio: child child is missing previous AddCLI")
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
