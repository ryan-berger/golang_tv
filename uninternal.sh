#!/usr/bin/env bash

#######################################
##### Un-internals your compiler ######
#######################################

# TODO, make this a bit better so go upgrades don't hurt really bad

mkdir internal/src
mkdir internal/src/cmd

cp -R go-linux-amd64-bootstrap/src/internal internal/src/internal
cp -R go-linux-amd64-bootstrap/src/cmd/compile/internal internal/src/cmd/compile
cp -R go-linux-amd64-bootstrap/src/cmd/internal internal/src/cmd/internal

find . -type f -name '*_test.go' -exec rm {} +

# remove abi_test.s in order to fix symbol error
rm internal/src/internal/abi/abi_test.s