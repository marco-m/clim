# 🫧  clim  🫧

[![Go Reference](https://pkg.go.dev/badge/github.com/marco-m/clim.svg)](https://pkg.go.dev/github.com/marco-m/clim)
[![Build Status](https://api.cirrus-ci.com/github/marco-m/clim.svg?branch=master)](https://cirrus-ci.com/github/marco-m/clim)

Command-line argument parsing for Go:

* Simple, small, can parse anything via the Value interface.
* No calls to os.Exit: easy to test, total control over termination.
* No output behind your back, always returns a string: easy to test.
* Support for subcommands.
* Support for help.
* Support for subcommand groups.

## Examples

See directory [examples](examples/).

## Status

Version 0.x, API can have breaking changes.

## Credits

Some code and inspiration taken from std/flag.
