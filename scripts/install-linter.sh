#!/bin/sh

set -ex

# From: https://golangci-lint.run/usage/install/#linux-and-windows
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v1.42.1
