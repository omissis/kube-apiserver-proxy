#!/bin/sh -x

set -e
set -o errexit -o nounset

go generate ./...
