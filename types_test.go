package clim_test

import (
	"testing"
	"time"

	"github.com/go-quicktest/qt"

	"github.com/marco-m/clim"
)

func TestParseIntSuccess(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want int
	}

	test := func(t *testing.T, tc testCase) {
		var count int
		cli := clim.New[any]("bang", "bangs head against wall", nil)
		cli.AddFlag(&clim.Flag{
			Value: clim.Int(&count, 3),
			Short: "c", Long: "count",
		})

		_, err := cli.Parse(tc.args)
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.Equals(count, tc.want))
	}

	testCases := []testCase{
		{
			name: "default value",
			args: nil,
			want: 3,
		},
		{
			name: "short",
			args: []string{"-c", "5"},
			want: 5,
		},
		{
			name: "long separated",
			args: []string{"--count", "7"},
			want: 7,
		},
		{
			name: "long with =",
			args: []string{"--count=9"},
			want: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestParseIntFailure(t *testing.T) {
	type testCase struct {
		name    string
		args    []string
		wantErr string
	}

	test := func(t *testing.T, tc testCase) {
		var count int
		cli := clim.New[any]("bang", "bangs head against wall", nil)
		cli.AddFlag(&clim.Flag{
			Value: clim.Int(&count, 3),
			Short: "c", Long: "count",
		})

		_, err := cli.Parse(tc.args)
		qt.Assert(t, qt.Equals(err.Error(), tc.wantErr))
	}

	testCases := []testCase{
		{
			name:    "not an int",
			args:    []string{"-c", "x"},
			wantErr: `setting "-c" "x": could not parse "x" as int (strconv.ParseInt: parsing "x": invalid syntax)`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestParseString(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want string
	}

	test := func(t *testing.T, tc testCase) {
		var fruit string
		cli := clim.New[any]("bang", "bangs head against wall", nil)
		cli.AddFlag(&clim.Flag{
			Value: clim.String(&fruit, "banana"),
			Short: "f", Long: "fruit",
		})

		_, err := cli.Parse(tc.args)
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.Equals(fruit, tc.want))
	}

	testCases := []testCase{
		{
			name: "default value",
			args: nil,
			want: "banana",
		},
		{
			name: "short",
			args: []string{"-f", "mango"},
			want: "mango",
		},
		{
			name: "long separated",
			args: []string{"--fruit", "tomato"},
			want: "tomato",
		},
		{
			name: "long with =",
			args: []string{"--fruit=papaya"},
			want: "papaya",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestParseBoolSuccess(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want bool
	}

	test := func(t *testing.T, tc testCase) {
		var sliced bool
		cli := clim.New[any]("bang", "bangs head against wall", nil)
		cli.AddFlag(&clim.Flag{
			Value: clim.Bool(&sliced, false),
			Short: "s", Long: "sliced",
		})

		_, err := cli.Parse(tc.args)
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.Equals(sliced, tc.want))
	}

	testCases := []testCase{
		{
			name: "default value",
			args: nil,
			want: false,
		},
		{
			name: "short",
			args: []string{"-s"},
			want: true,
		},
		{
			name: "long",
			args: []string{"--sliced"},
			want: true,
		},
		{
			name: "long with =",
			args: []string{"--sliced=true"},
			want: true,
		},
		{
			name: "long with =",
			args: []string{"--sliced=false"},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestParseBoolFailure(t *testing.T) {
	type testCase struct {
		name    string
		args    []string
		wantErr string
	}

	test := func(t *testing.T, tc testCase) {
		var sliced bool
		cli := clim.New[any]("bang", "bangs head against wall", nil)
		cli.AddFlag(&clim.Flag{
			Value: clim.Bool(&sliced, false),
			Short: "s", Long: "sliced",
		})

		_, err := cli.Parse(tc.args)
		qt.Assert(t, qt.Equals(err.Error(), tc.wantErr))
	}

	testCases := []testCase{
		{
			name:    "explicit value can be only true or false",
			args:    []string{"--sliced=ham"},
			wantErr: `setting "--sliced=ham": could not parse "ham" as bool (strconv.ParseBool: parsing "ham": invalid syntax)`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestParseDurationSuccess(t *testing.T) {
	var timeout time.Duration
	cli := clim.New[any]("bang", "bangs head against wall", nil)
	cli.AddFlag(&clim.Flag{
		Value: clim.Duration(&timeout, 0),
		Long:  "timeout",
	})

	_, err := cli.Parse([]string{"--timeout=32m4ms"})
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(timeout, 32*time.Minute+4*time.Millisecond))
}

func TestParseDurationFailure(t *testing.T) {
	var timeout time.Duration
	cli := clim.New[any]("bang", "bangs head against wall", nil)
	cli.AddFlag(&clim.Flag{
		Value: clim.Duration(&timeout, 0),
		Long:  "timeout",
	})

	_, err := cli.Parse([]string{"--timeout=78"})
	qt.Assert(t, qt.ErrorMatches(err,
		`setting "--timeout=78": time: missing unit in duration "78"`))
}
