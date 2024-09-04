# See https://just.systems/man/en/

# The first recipe is the one executed by default.
@_default:
    just --list

bang:
    go build -o ./bin/bang ./examples/bang

twocommands:
    go build -o ./bin/twocommands ./examples/twocommands

hg:
    go build -o ./bin/hg ./examples/hg

build:
    go build ./...

test:
    go test ./...
