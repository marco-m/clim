package clim

import "errors"

var (
	ErrHelp = errors.New("")
)

func Parse(cmdline []string, args any) error {
	return ErrHelp
}
