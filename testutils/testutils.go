package testutils

import (
	"io"
	"os"
	"testing"
)

// InterceptOutput replaces file (normally os.Stdout) with a pipe
// and returns the reset function. The reset function returns what
// has been written to file and resets it to point to the original
// (normally os.Stdout).
func InterceptOutput(t *testing.T, file **os.File) func() string {
	t.Helper()
	rd, wr, err := os.Pipe()
	if err != nil {
		t.Fatalf("replace: %s", err)
	}
	orig := *file
	*file = wr

	return func() string {
		t.Helper()
		wr.Close()
		out, err := io.ReadAll(rd)
		if err != nil {
			t.Fatalf("replace2: %s", err)
		}
		*file = orig
		return string(out)
	}
}
