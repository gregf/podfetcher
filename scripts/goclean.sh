#!/bin/bash
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. goimports     (https://github.com/bradfitz/goimports)
# 3. golint        (https://github.com/golang/lint)
# 4. go vet        (http://golang.org/cmd/vet)
# 5. test coverage (http://blog.golang.org/cover)

set -e

# Automatic checks
echo "gofmt..."
test -z "$(gofmt -l -w src | tee /dev/stderr)"
echo "goimports..."
test -z "$(goimports -l -w src | tee /dev/stderr)"
echo "golint..."
test -z "$(golint src | tee /dev/stderr)"
echo "go vet..."
go vet ./src/...
