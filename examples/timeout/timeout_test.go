package main

import (
	"testing"

	"github.com/go-quicktest/qt"
)

func TestTimeout(t *testing.T) {
	err := mainErr([]string{})
	qt.Assert(t, qt.ErrorMatches(err, `context deadline exceeded`))
}
