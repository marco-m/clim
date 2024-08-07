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

## How does it look like?

```console
$ go run ./examples/hg -h
hg -- Mercurial Distributed SCM

Usage: hg <command> [options]

available commands:

Repository creation:

 clone        make a copy of an existing repository
 init         create a new repository in the given directory

Remote repository management:

 incoming     show new changesets found in source
 outgoing     show changesets not found in the destination

Options:

 -h, --help    Print this help and exit
```

Let's follow a subcommand:

```console
$ go run ./examples/hg incoming -h
hg incoming -- show new changesets found in source

Usage: hg incoming [options]

Options:

 --bundle FILE             file to store the bundles into
 -f, --force               run even if remote repository is unrelated (default: false)
 -n, --newest-first        show newest record first (default: false)
 -r, --rev REV[,REV,..]    remote changeset(s) intended to be added

 -h, --help                Print this help and exit
```

## Examples

See directory [examples](examples/).

## Status

Version 0.x, API can have breaking changes.

## Credits

Some code and inspiration taken from std/flag.
