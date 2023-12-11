package clim_test

import (
	"testing"

	"github.com/go-quicktest/qt"

	"github.com/marco-m/clim"
)

func Test_SimpleHelp(t *testing.T) {
	type Args struct {
		Foo      string `clim:"argument required"`
		Bar      int    `clim:"argument"`
		Optimize int    `clim:"option"`
	}
	var args Args
	want := `
Usage:
  example [--optimize <value>] <foo> [<bar>]

Arguments:
  <foo>
  <bar>

Options:
  --optimize <value>
  --help, -h          display this help and exit`

	err := clim.Parse([]string{"-h"}, &args)

	qt.Assert(t, qt.ErrorIs(err, clim.ErrHelp))
	qt.Assert(t, qt.Equals(err.Error(), want))
}

func Test_SimpleSuccess(t *testing.T) {
	type Args struct {
		Foo string
		Bar int
	}

	type testCase struct {
		name    string
		cmdline []string
		want    Args
	}

	test := func(t *testing.T, tc testCase) {
		var args Args
		err := clim.Parse(tc.cmdline, &args)

		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.DeepEquals(args, tc.want))
	}

	testCases := []testCase{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}

func Test_SimpleFailure(t *testing.T) {
}
