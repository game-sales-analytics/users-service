#!/bin/sh

set -ex

rm -vrf ./bin
mkdir -vp ./bin
PKGS=$(go list -find -f '{{.ImportPath}}' ./cmd/...)
export PKGS
for pkg in $PKGS; do
  echo "building package $pkg... "
  go build -o ./bin/ "$pkg"
  echo 'Done.'
done
