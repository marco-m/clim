package clim_test

import (
	"testing"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestPosArgsRequiredSuccess(t *testing.T) {
	cli, err := clim.New[any](nil, "bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	var positionals []string
	err = cli.AddPosArgs(&positionals,
		clim.Pair{"NAME", "Name of the foos (required)"},
		clim.Pair{"COUNT", "How many foos (required)"})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"mangos", "7"})
	rosina.AssertNoError(t, err)
	// Yes, the first implementation is almost cheating.
	want := []string{"mangos", "7"}
	rosina.AssertDeepEqual(t, positionals, want, "positionals")
}

func TestCannotAddPosArgAfterSubCommand(t *testing.T) {
	cli, err := clim.New[any](nil, "bang", "bang head", nil)
	rosina.AssertNoError(t, err)

	_, err = clim.New[any](cli, "sub", "I am a subcommand A", nil)
	rosina.AssertNoError(t, err)

	var positionals []string
	err = cli.AddPosArgs(&positionals, clim.Pair{"NAME", "Name of the foos"})
	rosina.AssertErrorContains(t, err,
		"bang: already have subcommands; cannot have also pos args")
}

func TestAddPosArgFailure(t *testing.T) {
	type testCase struct {
		name  string
		pairs []clim.Pair
		want  string
	}

	test := func(t *testing.T, tc testCase) {
		cli, err := clim.New[any](nil, "bang", "bang head", nil)
		rosina.AssertNoError(t, err)

		var positionals []string
		err = cli.AddPosArgs(&positionals, tc.pairs...)
		rosina.AssertErrorContains(t, err, tc.want)
	}

	testCases := []testCase{
		{
			name:  "already defined",
			pairs: []clim.Pair{{"A", "foo"}, {"A", "bar"}},
			want:  `bang: pos arg at index 1 ("A") was already defined at index 0`,
		},
		{
			name:  "empty name",
			pairs: []clim.Pair{{"", "foo"}},
			want:  `bang: pos arg at index 0 ("") cannot be empty`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

// func TestPosArgsOptional(t *testing.T) {
// 	want := `bang -- bang head

// Usage: bang [options] [NAME [COUNT]]

// Options:

//  -h, --help    Print this help and exit

// Positional arguments:

//  NAME       Name of the foos (default: banana)
//  COUNT      How many foos (default: 42)
// `
// }

// Other test:
// // No subcommand nor positional.
// // Ensure that nothing is remaining unprocessed.
// if len(cli.positionals) > 0 {
// 	return nil, NewParseError("%s: unrecognized arguments: %s",
// 		cli.rootToHere, strings.Join(cli.positionals, " "))
// }

// func TestPosArgsSimpleSuccess(t *testing.T) {
// 	// assing and check assignment
// 	t.Fatal("writeme")
// }
