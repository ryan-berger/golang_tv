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

### Example:
```go
func Src(a, b int) int {
	return (a >> b) << b
}

func Tgt(a, b int) int {
	return a >> (b << b)
}
```

We feed it to golang_tv and it spits out:
```
Src func(int, int) int // Go's SSA
  b1:
    (?) v1 = InitMem <mem>
    (-3) v7 = ArgIntReg <int> {a+0} [0] (a[int])
    (-3) v8 = ArgIntReg <int> {b+0} [1] (b[int])
    (?) v9 = Const64 <int> [0]
    (+4) v10 = Leq64 <bool> v9 v8
    (?) v11 = Const64 <int64> [0] DEAD
    If v10 -> b4 b3 (likely)
  b2: DEAD
    BlockInvalid
  b3: <- b1
    (4) v12 = StaticCall <mem> {AuxCall{runtime.panicshift}} v1
    Exit v12
  b4: <- b1
    (4) v14 = Rsh64x64 <int> [false] v7 v8
    (+4) v17 = Lsh64x64 <int> [false] v14 v8
    (4) v18 = MakeResult <int,mem> v17 v1
    Ret v18
name a[int]: [v7]
name b[int]: [v8]
 Tgt func(int, int) int
  b1:
    (?) v1 = InitMem <mem>
    (-7) v7 = ArgIntReg <int> {a+0} [0] (a[int])
    (-7) v8 = ArgIntReg <int> {b+0} [1] (b[int])
    (?) v9 = Const64 <int> [0]
    (+8) v10 = Leq64 <bool> v9 v8
    (?) v6 = Const64 <int64> [0] DEAD
    If v10 -> b2 b3 (likely)
  b2: <- b1
    (8) v14 = Lsh64x64 <int> [false] v8 v8
    (8) v15 = Leq64 <bool> v9 v14
    If v15 -> b4 b3 (likely)
  b3: <- b1 b2
    (8) v12 = StaticCall <mem> {AuxCall{runtime.panicshift}} v1
    Exit v12
  b4: <- b2
    (8) v16 = Rsh64x64 <int> [false] v7 v14
    (8) v17 = MakeResult <int,mem> v16 v1
    Ret v17
name a[int]: [v7]
name b[int]: [v8]


---------------------------------------- // Alive IR
define {i64} @src({i64, i64} %0) {
%b1:
  %v7 = extractvalue {i64, i64} %0, 0
  %v8 = extractvalue {i64, i64} %0, 1
  %v10 = icmp sle i64 0, %v8
  br i1 %v10, label %b4, label %b3

%b3:
  assume i1 0

%b4:
  %should_clamp = icmp sge i64 %v8, 63
  %clamped = select i1 %should_clamp, i64 63, i64 %v8
  %v14 = ashr i64 %v7, %clamped
  %should_clamp1 = icmp sge i64 %v8, 63
  %clamped2 = select i1 %should_clamp1, i64 63, i64 %v8
  %v17 = lshr i64 %v14, %clamped2
  %first = insertvalue {i64} undef, i64 %v17, 0
  ret {i64} %first
}
=>
define {i64} @tgt({i64, i64} %0) {
%b1:
  %v7 = extractvalue {i64, i64} %0, 0
  %v8 = extractvalue {i64, i64} %0, 1
  %v10 = icmp sle i64 0, %v8
  br i1 %v10, label %b2, label %b3

%b2:
  %should_clamp = icmp sge i64 %v8, 63
  %clamped = select i1 %should_clamp, i64 63, i64 %v8
  %v14 = lshr i64 %v8, %clamped
  %v15 = icmp sle i64 0, %v14
  br i1 %v15, label %b4, label %b3

%b4:
  %should_clamp1 = icmp sge i64 %v14, 63
  %clamped2 = select i1 %should_clamp1, i64 63, i64 %v14
  %v16 = ashr i64 %v7, %clamped2
  %first = insertvalue {i64} undef, i64 %v16, 0
  ret {i64} %first

%b3:
  assume i1 0
}
Transformation doesn't verify!

ERROR: Value mismatch

Example:
{i64, i64} %0 = { #x0000000000000004 (4), #x0000000000000020 (32) }

Source:
i64 %v7 = #x0000000000000004 (4)
i64 %v8 = #x0000000000000020 (32)
i1 %v10 = #x1 (1)
  >> Jump to %b4
i1 %should_clamp = #x0 (0)
i64 %clamped = #x0000000000000020 (32)
i64 %v14 = #x0000000000000000 (0)
i1 %should_clamp1 = #x0 (0)
i64 %clamped2 = #x0000000000000020 (32)
i64 %v17 = #x0000000000000000 (0)
{i64} %first = { #x0000000000000000 (0) }

Target:
i64 %v7 = #x0000000000000004 (4)
i64 %v8 = #x0000000000000020 (32)
i1 %v10 = #x1 (1)
  >> Jump to %b2
i1 %should_clamp = #x0 (0)
i64 %clamped = #x0000000000000020 (32)
i64 %v14 = #x0000000000000000 (0)
i1 %v15 = #x1 (1)
  >> Jump to %b4
i1 %should_clamp1 = #x0 (0)
i64 %clamped2 = #x0000000000000000 (0)
i64 %v16 = #x0000000000000004 (4)
{i64} %first = { #x0000000000000004 (4) }
Source value: { #x0000000000000000 (0) }
Target value: { #x0000000000000004 (4) }
```
Looks like `(a >> b) << b` is not the same as `a >> (b << b)`! The example it gives is when a=4, b=32, which causes Src to retur 0 and Tgt to return 4

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
