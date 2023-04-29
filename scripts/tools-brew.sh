#!/bin/sh -x

set -e
set -o errexit -o nounset

brew install checkmake markdownlint-cli2 jsonlint
