#!/bin/sh -x

set -e
set -o errexit -o nounset

go install github.com/daixiang0/gci@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/nikolaydubina/go-cover-treemap@latest
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
