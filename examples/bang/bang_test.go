package main

import (
	"os"
	"testing"

	"github.com/marco-m/rosina"
)

func TestBang(t *testing.T) {
	want := `1 bang against cardboard
2 bang against cardboard
3 bang against cardboard
`
	reset := rosina.InterceptOutput(t, &os.Stdout)

	err := mainErr([]string{})
	rosina.AssertNoError(t, err)

	out := reset()
	rosina.AssertEqual(t, out, want, "stdout")
}
