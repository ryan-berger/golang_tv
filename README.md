# golang_tv

golang_tv is a translation validation platform built upon Alive2

The entire project is in its infancy, so many, many features are missing.

## Conceptual Overview

golang_tv exposes Go compiler internals via the `cmd/compile/internal/ssa`
package in order to lift to LLVM, and verify via Alive2.

It makes use of the fact that LLVM IR is in quasi-SSA form, and Go SSA is in SSA form
to make easy translation of Go SSA's generic ops to LLVM IR.

It is not a direct translation however, because LLVM IR must simulate
Go semantics to get use out of it.

## Setup

The setup is extremely tedious, and is not particularly automated. This was a hack I put together in 3 days.
If you want better in the immediate future, PRs are welcome, I am actively working on making it better.

You first will need to build **DEBUG BUILDS** LLVM and AliveToolkit/Alive2 from source with clang (use gcc/g++ at your own risk). 

You also should grab the source of ryan-berger/go-llvm and checkout the golang-tv branch.
You will likely have to edit the linking commands to fit your install paths, so prepare to use a `go.work`
to link things up.

Alive2 and tinygo-org/go-llvm have some different needs from the LLVM build,
so instead of running Alive's LLVM build command of:
```
cd llvm
mkdir build
cd build
cmake -GNinja -DLLVM_ENABLE_RTTI=ON -DLLVM_ENABLE_EH=ON -DBUILD_SHARED_LIBS=ON -DCMAKE_BUILD_TYPE=Debug -DLLVM_TARGETS_TO_BUILD=X86 -DLLVM_ENABLE_ASSERTIONS=ON -DLLVM_ENABLE_PROJECTS="llvm;clang" ../llvm
```

Run:
```
cd llvm
mkdir build
cd build
cmake -GNinja -DLLVM_ENABLE_RTTI=ON -DLLVM_ENABLE_EH=ON -DBUILD_SHARED_LIBS=ON -DCMAKE_BUILD_TYPE=Debug -DLLVM_ENABLE_ASSERTIONS=ON -DLLVM_ENABLE_PROJECTS="llvm;clang" ../llvm
```

Change the end of Alive2's build command from 

```
Release ..
```

to 
```
Debug ..
```

So that tinygo-org/go-llvm can get all the backend generation dependencies it needs.

With this all set up, you should be good to go.

## Validation

golang_tv only supports validation between two Go functions with the names of
Source and Target, which are read in from `tests/validate.go`.

In the future, it should also support one function that validates 
that all Go optimization passes are valid

## Contributing

A simple PR will do!

### Things to contribute on:
- There are a ton more instructions to support in `internal/gen/ssa2llvm.go`
and help would be very helpful.

- Better Go ABI support would be helpful to better model the stack-based calling convention
that Go uses. Currently only register based calls are supported, anything else will be weird 

- The build system isn't particularly great and is very hardcoded. Automated setup would be awesome

- uninternal.sh needs to download new go builds and pull out compiler pieces
into internal/src like we have now

- The SSA API that I have exposed is extremely hacky because it was never meant to
be exposed to anyone outside the Go compiler team. Clean up and clever hooks into the API,
documentation, etc would be extremely helpful.

- get `panic` working