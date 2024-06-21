# See https://just.systems/man/en/

# The first recipe is the one executed by default.
@_default:
    just --list

build:
    go build ./...

test:
    go test ./...
