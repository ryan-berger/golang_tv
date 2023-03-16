#!/usr/bin/env bash

mkdir internal/src
mkdir internal/src/cmd

cp -R go-linux-amd64-bootstrap/src/internal internal/src/internal
cp -R go-linux-amd64-bootstrap/src/cmd/compile/internal internal/src/cmd/compile
cp -R go-linux-amd64-bootstrap/src/cmd/internal internal/src/cmd/internal

find . -type f -name '*_test.go' -exec rm {} +

# remove abi_test.s in order to fix symbol error
rm internal/src/internal/abi/abi_test.s