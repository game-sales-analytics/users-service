#!/bin/sh

set -ex

export LINTER=${GOBIN:-"$GOPATH/bin"}/golangci-lint

ensure_linter_is_installed() {
  command -v $LINTER && return
  echo 'golangci-lint command not found! installing...'
  # sourcing other script with relative path
  # source: https://stackoverflow.com/a/48167380
  . "$(dirname $0)/install-linter.sh"
}

ensure_linter_is_installed

$LINTER run \
  --print-resources-usage \
  --sort-results \
  --verbose \
  --no-config \
  --print-linter-name \
  --skip-files ".*.gen.go" \
  --sort-results \
  --deadline 7m \
  --tests \
  --enable=structcheck \
  --enable=deadcode \
  --enable=gocyclo \
  --enable=ineffassign \
  --enable=revive \
  --enable=goimports \
  --enable=errcheck \
  --enable=varcheck \
  --enable=goconst \
  --enable=megacheck \
  --enable=misspell \
  --enable=unused \
  --enable=typecheck \
  --enable=staticcheck \
  --enable=govet \
  --enable=gosimple \
  ./...
