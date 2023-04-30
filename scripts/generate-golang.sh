#!/bin/sh -x

set -e
set -o errexit -o nounset

go run github.com/99designs/gqlgen generate && go generate ./...
