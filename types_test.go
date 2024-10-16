package clim_test

import (
	"testing"
	"time"

	"github.com/marco-m/clim"
	"github.com/marco-m/rosina"
)

func TestParseIntSuccess(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want int
	}

	test := func(t *testing.T, tc testCase) {
		var count int
		cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
		rosina.AssertNoError(t, err)

		err = cli.AddFlags(&clim.Flag{
			Value: clim.Int(&count, 3),
			Short: "c", Long: "count",
		})
		rosina.AssertNoError(t, err)

		_, err = cli.Parse(tc.args)
		rosina.AssertNoError(t, err)
		rosina.AssertEqual(t, count, tc.want, "count")
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
		cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
		rosina.AssertNoError(t, err)

		err = cli.AddFlags(&clim.Flag{
			Value: clim.Int(&count, 3),
			Short: "c", Long: "count",
		})
		rosina.AssertNoError(t, err)

		_, err = cli.Parse(tc.args)
		rosina.AssertErrorIs(t, err, clim.ErrParse)
		rosina.AssertErrorContains(t, err, tc.wantErr)
	}

	testCases := []testCase{
		{
			name:    "not an int",
			args:    []string{"-c", "x"},
			wantErr: `setting "-c" "x": could not parse "x" as int (strconv.ParseInt: parsing "x": invalid syntax)`,
		},
		{
			name:    "missing value",
			args:    []string{"-c"},
			wantErr: `flag "-c" requires a value`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { test(t, tc) })
	}
}

func TestParseIntSliceSuccess(t *testing.T) {
	var pippos []int
	cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value: clim.IntSlice(&pippos, []int{10}),
		Short: "p", Long: "pippos",
	})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--pippos=1,2,3"})
	rosina.AssertNoError(t, err)
	rosina.AssertDeepEqual(t, pippos, []int{1, 2, 3}, "pippos")
}

func TestParseIntSliceFailure(t *testing.T) {
	var pippos []int
	cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value: clim.IntSlice(&pippos, nil),
		Short: "p", Long: "pippos",
	})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--pippos=a,b,c"})
	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err,
		`setting "--pippos=a,b,c": could not parse "a" as int (strconv.Atoi: parsing "a": invalid syntax)`)
}

func TestParseString(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want string
	}

	test := func(t *testing.T, tc testCase) {
		var fruit string
		cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
		rosina.AssertNoError(t, err)

		err = cli.AddFlags(&clim.Flag{
			Value: clim.String(&fruit, "banana"),
			Short: "f", Long: "fruit",
		})
		rosina.AssertNoError(t, err)

		_, err = cli.Parse(tc.args)
		rosina.AssertNoError(t, err)
		rosina.AssertEqual(t, fruit, tc.want, "fruit")
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

func TestParseStringSliceSuccess(t *testing.T) {
	var mickeys []string
	cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value: clim.StringSlice(&mickeys, []string{"x"}),
		Short: "m", Long: "mickeys",
	})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--mickeys=a,b,c"})
	rosina.AssertNoError(t, err)
	rosina.AssertDeepEqual(t, mickeys, []string{"a", "b", "c"}, "mickeys")
}

func TestParseBoolSuccess(t *testing.T) {
	type testCase struct {
		name string
		args []string
		want bool
	}

	test := func(t *testing.T, tc testCase) {
		var sliced bool
		cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
		rosina.AssertNoError(t, err)

		err = cli.AddFlags(&clim.Flag{
			Value: clim.Bool(&sliced, false),
			Short: "s", Long: "sliced",
		})
		rosina.AssertNoError(t, err)

		_, err = cli.Parse(tc.args)
		rosina.AssertNoError(t, err)
		rosina.AssertEqual(t, sliced, tc.want, "sliced")
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
			name: "explicit value, true",
			args: []string{"--sliced=true"},
			want: true,
		},
		{
			name: "explicit value, false",
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
		cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
		rosina.AssertNoError(t, err)

		err = cli.AddFlags(&clim.Flag{
			Value: clim.Bool(&sliced, false),
			Short: "s", Long: "sliced",
		})
		rosina.AssertNoError(t, err)

		_, err = cli.Parse(tc.args)
		rosina.AssertErrorIs(t, err, clim.ErrParse)
		rosina.AssertErrorContains(t, err, tc.wantErr)
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
	cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value: clim.Duration(&timeout, 0),
		Long:  "timeout",
	})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--timeout=32m4ms"})
	rosina.AssertNoError(t, err)
	rosina.AssertEqual(t, timeout, 32*time.Minute+4*time.Millisecond, "timeout")
}

func TestParseDurationFailure(t *testing.T) {
	var timeout time.Duration
	cli, err := clim.NewTop[any]("bang", "bangs head against wall", nil)
	rosina.AssertNoError(t, err)

	err = cli.AddFlags(&clim.Flag{
		Value: clim.Duration(&timeout, 0),
		Long:  "timeout",
	})
	rosina.AssertNoError(t, err)

	_, err = cli.Parse([]string{"--timeout=78"})
	rosina.AssertErrorIs(t, err, clim.ErrParse)
	rosina.AssertErrorContains(t, err,
		`setting "--timeout=78": time: missing unit in duration "78"`)
}
